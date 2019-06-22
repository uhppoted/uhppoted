package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"uhppote"
	codec "uhppote/encoding/UTO311-L0x"
)

var handlers = map[byte]func(*net.UDPConn, *net.UDPAddr, []byte){
	0x94: find,
	0x50: putCard,
	0x52: deleteCard,
	0x58: getCards,
	0x5a: getCardById,
	0x5c: getCardByIndex,
}

func simulate() {
	interrupt := make(chan int, 1)
	bind, err := net.ResolveUDPAddr("udp", ":60000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP bind address [%v]", err)))
		return
	}

	go run(bind, interrupt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	stop(interrupt)

	os.Exit(1)
}

func stop(interrupt chan int) {
	interrupt <- 42
	<-interrupt
}

func run(bind *net.UDPAddr, interrupt chan int) {
	connection, err := net.ListenUDP("udp", bind)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to bind to UDP socket [%v]", err)))
		return
	}

	defer connection.Close()

	wait := make(chan int, 1)
	go func() {
		err := listenAndServe(connection)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0
	}()

	<-interrupt
	connection.Close()
	<-wait
	interrupt <- 0
}

func listenAndServe(c *net.UDPConn) error {
	for {
		request, remote, err := receive(c)
		if err != nil {
			return err
		}

		handle(c, remote, request)
	}

	return nil
}

func handle(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	if len(bytes) != 64 {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid message length %d", len(bytes))))
		return
	}

	if bytes[0] != 0x17 {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid message type %02X", bytes[0])))
		return
	}

	h := handlers[bytes[1]]
	if h == nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid command %02X", bytes[1])))
		return
	}

	h(c, src, bytes)
}

func receive(c *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	request := make([]byte, 2048)

	N, remote, err := c.ReadFromUDP(request)
	if err != nil {
		return []byte{}, nil, errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if options.debug {
		fmt.Printf(" ... received %v bytes from %v\n%s\n", N, remote, dump(request[0:N], " ...          "))
	}

	return request[:N], remote, nil
}

func send(c *net.UDPConn, dest *net.UDPAddr, response interface{}) {
	if response == nil {
		return
	}

	message, err := codec.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	N, err := c.WriteTo(message, dest)
	if err != nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err)))
	} else {
		if options.debug {
			fmt.Printf(" ... sent %v bytes to %v\n%s\n", N, dest, dump(message[0:N], " ...          "))
		}
	}
}

func find(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.FindDevicesRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response, err := s.Find(bytes)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			send(c, src, response)
		}
	}
}

func putCard(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.PutCardRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		if s.SerialNumber == request.SerialNumber {
			response, err := s.PutCard(request)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func deleteCard(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.DeleteCardRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		if s.SerialNumber == request.SerialNumber {
			response, err := s.DeleteCard(request)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func getCards(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.GetCardsRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		if s.SerialNumber == request.SerialNumber {
			response, err := s.GetCards(request)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func getCardById(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.GetCardByIdRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response, err := s.GetCardById(request)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			send(c, src, response)
		}
	}
}

func getCardByIndex(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request := uhppote.GetCardByIndexRequest{}

	err := codec.Unmarshal(bytes, &request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response, err := s.GetCardByIndex(request)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			send(c, src, response)
		}
	}
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
