package uhppote

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
)

type UHPPOTE struct {
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
	Debug            bool
}

func (u *UHPPOTE) Execute(request, reply interface{}) error {
	bind, err := u.bindAddr()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to resolve bind address [%v]", err))
	}

	dest, err := u.broadcastAddr()
	if err != nil {
		return err
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

func (u *UHPPOTE) Broadcast(request interface{}) ([][]byte, error) {
	addr, err := u.broadcastAddr()
	if err != nil {
		return [][]byte{}, errors.New(fmt.Sprintf("Unable to resolve UDP broadcast address [%v]", err))
	}

	return u.broadcast(request, addr)
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

	bindTo, err := u.bindAddr()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to resolve bind address [%v]", err))
	}

	connection, err := net.ListenUDP("udp", bindTo)
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
	if u.BindAddress == nil {
		return errors.New("Listen requires a fixed bind address")
	}

	c, err := u.open(u.BindAddress)
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

func localAddr() (*net.UDPAddr, error) {
	list, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range list {
		addrs, err := iface.Addrs()
		if err == nil {
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err == nil && !ip.IsLoopback() && ip.To4() != nil {
					return &net.UDPAddr{ip.To4(), 0, ""}, nil
				}
			}
		}
	}

	return nil, errors.New("Unable to identify local interface to bind to")
}

func (u *UHPPOTE) bindAddr() (*net.UDPAddr, error) {
	if u.BindAddress != nil {
		return u.BindAddress, nil
	}

	addr, err := localAddr()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to resolve bind address [%v]", err))
	}

	return addr, nil
}

// Ref. https://stackoverflow.com/questions/36166791/how-to-get-broadcast-address-of-ipv4-net-ipnet
func (u *UHPPOTE) broadcastAddr() (*net.UDPAddr, error) {
	if u.BroadcastAddress != nil {
		return u.BroadcastAddress, nil
	}

	if u.BindAddress != nil && u.BindAddress.IP.IsUnspecified() {
		return net.ResolveUDPAddr("udp", "255.255.255.255:60000")
	}

	local, err := u.bindAddr()
	if err != nil {
		return nil, err
	}

	if local.IP.To4() == nil {
		return nil, errors.New("IPv6 does not support broadcast addresses")
	}

	broadcast := make(net.IP, len(local.IP.To4()))
	addr := binary.BigEndian.Uint32(local.IP.To4()) | ^binary.BigEndian.Uint32(net.IP(local.IP.DefaultMask()).To4())

	binary.BigEndian.PutUint32(broadcast, addr)

	return &net.UDPAddr{broadcast.To4(), 60000, ""}, nil
}
