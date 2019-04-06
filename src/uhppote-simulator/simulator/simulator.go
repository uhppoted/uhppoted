package simulator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"time"
)

type Simulator struct {
	Interrupt chan int
	Debug     bool
}

func (s *Simulator) Stop() {
	s.Interrupt <- 42
	<-s.Interrupt
}

func (s *Simulator) Run() {
	local, err := net.ResolveUDPAddr("udp", "0.0.0.0:60000")

	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP local address [%v]", err)))
		return
	}

	connection, err := net.ListenUDP("udp", local)

	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err)))
		return
	}

	defer connection.Close()

	wait := make(chan int, 1)
	go func() {
		err := s.listenAndServe(connection)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0

	}()

	<-s.Interrupt
	connection.Close()
	<-wait
	s.Interrupt <- 0
}

func (s *Simulator) listenAndServe(c *net.UDPConn) error {
	request := make([]byte, 2048)

	for {
		N, remote, err := c.ReadFromUDP(request)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
		}

		if s.Debug {
			regex := regexp.MustCompile("(?m)^(.*)")

			fmt.Printf(" ... received %v bytes from %v\n", N, remote)
			fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(request[:N]), " ... $1"))
		}

		// ...parse request

		reply, err := s.handle(request[:N])

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			N, err = c.WriteTo(reply, remote)

			if err != nil {
				return errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
			}

			if s.Debug {
				regex := regexp.MustCompile("(?m)^(.*)")

				fmt.Printf(" ... sent %v bytes to %v\n", N, remote)
				fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ... $1"))
			}
		}
	}

	return nil
}

func (s *Simulator) handle(bytes []byte) ([]byte, error) {
	if len(bytes) != 64 {
		return []byte{}, errors.New(fmt.Sprintf("Invalid message length %d", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return []byte{}, errors.New(fmt.Sprintf("Invalid message type %02X", bytes[0]))
	}

	switch bytes[1] {
	case 0x94:
		return s.search(bytes)
	default:
		return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
	}

	return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
}

func parse(bytes []byte) interface{} {
	fmt.Printf("%v %x %x\n", len(bytes), bytes[0], bytes[1])
	if len(bytes) == 64 && bytes[0] == 0x17 {
		switch bytes[1] {
		case 0x94:
			request := struct {
				MsgType byte `uhppote:"offset:1"`
			}{
				0x94,
			}

			return &request
		}
	}

	return nil
}

func (s *Simulator) search(bytes []byte) ([]byte, error) {
	time.Sleep(100 * time.Millisecond)

	msg := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x01, 0x7d, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	return msg, nil
}
