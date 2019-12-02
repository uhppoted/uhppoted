package uhppote

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
)

var VERSION string = "v0.5.0"

type UHPPOTE struct {
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
	Devices          map[uint32]*net.UDPAddr
	Debug            bool
}

func (u *UHPPOTE) Send(serialNumber uint32, request interface{}) (messages.Response, error) {
	bind := u.bindAddress()
	dest := u.Devices[serialNumber]

	if dest == nil {
		dest = u.broadcastAddress()
	}

	c, err := u.open(bind)
	if err != nil {
		return nil, err
	}

	defer func() {
		c.Close()
	}()

	if err = u.send(c, dest, request); err != nil {
		return nil, err
	}

	err = c.SetDeadline(time.Now().Add(5000 * time.Millisecond))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to set UDP timeout [%v]", err))
	}

	m := make([]byte, 2048)
	N, remote, err := c.ReadFromUDP(m)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... received %v bytes from %v\n ... response\n%s\n", N, remote, dump(m[:N], " ...          "))
	}

	response, err := messages.UnmarshalResponse(m[:N])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (u *UHPPOTE) Execute(serialNumber uint32, request, reply interface{}) error {
	bind := u.bindAddress()
	dest := u.Devices[serialNumber]

	if dest == nil {
		dest = u.broadcastAddress()
	}

	c, err := u.open(bind)
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
	}()

	if err = u.send(c, dest, request); err == nil {
		if reply != nil {
			return u.receive(c, reply)
		}
	}

	return err
}

func (u *UHPPOTE) Broadcast(request, replies interface{}) error {
	if m, err := u.broadcast(request, u.broadcastAddress()); err != nil {
		return err
	} else {
		return codec.UnmarshalArray(m, replies)
	}
}

// Sends a UDP message to a specific device but anticipates replies from more than one device because
// it may fall back to the broadcast address if the device ID has no configured IP address.
func (u *UHPPOTE) DirectedBroadcast(serialNumber uint32, request, replies interface{}) error {
	dest := u.Devices[serialNumber]
	if dest == nil {
		dest = u.broadcastAddress()
	}

	if m, err := u.broadcast(request, dest); err != nil {
		return err
	} else {
		return codec.UnmarshalArray(m, replies)
	}
}

func (u *UHPPOTE) open(addr *net.UDPAddr) (*net.UDPConn, error) {
	connection, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err))
	}

	return connection, nil
}

func (u *UHPPOTE) send(connection *net.UDPConn, addr *net.UDPAddr, request interface{}) error {
	m, err := codec.Marshal(request)
	if err != nil {
		return err
	}

	if u.Debug {
		fmt.Printf(" ... request\n%s\n", dump(m, " ...          "))
	}

	N, err := connection.WriteTo(m, addr)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes to %v\n", N, addr)
	}

	return nil
}

func (u *UHPPOTE) broadcast(request interface{}, addr *net.UDPAddr) ([][]byte, error) {
	m, err := codec.Marshal(request)
	if err != nil {
		return nil, err
	}

	if u.Debug {
		fmt.Printf(" ... request\n%s\n", dump(m, " ...          "))
	}

	bind := u.bindAddress()
	connection, err := net.ListenUDP("udp", bind)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err))
	}

	defer func() {
		connection.Close()
	}()

	N, err := connection.WriteTo(m, addr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes to %v\n", N, addr)
	}

	replies := make([][]byte, 0)
	go func() {
		for {
			reply := make([]byte, 2048)
			N, remote, err := connection.ReadFromUDP(reply)

			if err != nil {
				break
			} else {
				replies = append(replies, reply[:N])

				if u.Debug {
					regex := regexp.MustCompile("(?m)^(.*)")

					fmt.Printf(" ... received %v bytes from %v\n", N, remote)
					fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ...          $1"))
				}
			}
		}
	}()

	time.Sleep(2500 * time.Millisecond)

	return replies, err
}

func (u *UHPPOTE) receive(c *net.UDPConn, reply interface{}) error {
	m := make([]byte, 2048)

	err := c.SetDeadline(time.Now().Add(5000 * time.Millisecond))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to set UDP timeout [%v]", err))
	}

	N, remote, err := c.ReadFromUDP(m)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... received %v bytes from %v\n ... response\n%s\n", N, remote, dump(m[:N], " ...          "))
	}

	return codec.Unmarshal(m[:N], reply)
}

func (u *UHPPOTE) listen(p chan Event, q chan os.Signal) error {
	bind := u.bindAddress()
	if bind.Port == 0 {
		return errors.New("Listen requires a non-zero UDP port")
	}

	c, err := u.open(bind)
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
	}()

	closed := false
	go func() {
		for {
			if s := <-q; s == os.Interrupt {
				closed = true
				c.Close()
			}
		}
	}()

	m := make([]byte, 2048)

	for {
		if u.Debug {
			fmt.Printf(" ... listening\n")
		}

		N, remote, err := c.ReadFromUDP(m)
		if err != nil {
			if closed {
				return nil
			}

			return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
		}

		if u.Debug {
			fmt.Printf(" ... received %v bytes from %v\n ... response\n%s\n", N, remote, dump(m[:N], " ...          "))
		}

		event := Event{}
		err = codec.Unmarshal(m[:N], &event)
		if err != nil {
			return errors.New(fmt.Sprintf("Error unmarshalling event [%v]", err))
		}

		p <- event
	}

	return nil
}

func (u *UHPPOTE) bindAddress() *net.UDPAddr {
	if u.BindAddress != nil {
		return u.BindAddress
	}

	addr := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 0,
		Zone: "",
	}

	copy(addr.IP, net.IPv4zero)

	return &addr
}

func (u *UHPPOTE) broadcastAddress() *net.UDPAddr {
	if u.BroadcastAddress != nil {
		return u.BroadcastAddress
	}

	addr := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	copy(addr.IP, net.IPv4bcast)

	return &addr
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
