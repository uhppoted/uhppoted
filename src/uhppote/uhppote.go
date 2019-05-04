package uhppote

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
)

type UHPPOTE struct {
	BindAddress net.UDPAddr
	Debug       bool
}

func (u *UHPPOTE) Execute(request, reply interface{}) error {
	c, err := u.open()
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
	}()

	if err = u.send(c, request); err == nil {
		if reply != nil {
			return u.receive(c, reply)
		}
	}

	return err
}

func (u *UHPPOTE) Broadcast(request interface{}) ([][]byte, error) {
	p, err := codec.Marshal(request)
	if err != nil {
		return [][]byte{}, err
	}

	return u.broadcast(p)
}

func (u *UHPPOTE) listen(event interface{}) error {
	c, err := u.open()
	if err != nil {
		return err
	}

	defer func() {
		c.Close()
	}()

	m := make([]byte, 2048)
	N, remote, err := c.ReadFromUDP(m)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... received %v bytes from %v\n ... response\n%s\n", N, remote, dump(m[:N], " ...          "))
	}

	return codec.Unmarshal(m[:N], event)
}

func (u *UHPPOTE) open() (*net.UDPConn, error) {
	connection, err := net.ListenUDP("udp", &u.BindAddress)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err))
	}

	return connection, nil
}

func (u *UHPPOTE) send(connection *net.UDPConn, request interface{}) error {
	m, err := codec.Marshal(request)
	if err != nil {
		return err
	}

	if u.Debug {
		fmt.Printf(" ... request\n%s\n", dump(m, " ...          "))
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")
	//broadcast, err := net.ResolveUDPAddr("udp", "192.168.1.255:60000")

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to resolve UDP broadcast address [%v]", err))
	}

	N, err := connection.WriteTo(m, broadcast)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

	return nil
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

func (u *UHPPOTE) broadcast(cmd []byte) ([][]byte, error) {
	replies := make([][]byte, 0)

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... command %v bytes\n", len(cmd))
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(cmd), " ...         $1"))
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to resolve UDP broadcast address [%v]", err))
	}

	connection, err := net.ListenUDP("udp", &u.BindAddress)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err))
	}

	defer func() {
		connection.Close()
	}()

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
	}

	if u.Debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

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
