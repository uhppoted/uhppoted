package simulator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"regexp"
	"time"
	"uhppote"
	codec "uhppote/encoding"
	"uhppote/types"
)

type Simulator struct {
	Interrupt    chan int
	Debug        bool
	SerialNumber types.SerialNumber
	IpAddress    net.IP
	SubnetMask   net.IP
	Gateway      net.IP
	MacAddress   net.HardwareAddr
	Version      types.Version
	Date         types.Date
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
		return s.find(bytes)
	default:
		return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
	}

	return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
}

func (s *Simulator) find(bytes []byte) ([]byte, error) {
	time.Sleep(100 * time.Millisecond)

	response := uhppote.FindDevicesResponse{
		MsgType:      0x94,
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   s.MacAddress,
		Version:      s.Version,
		Date:         s.Date,
	}

	reply, err := codec.Marshal(response)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
