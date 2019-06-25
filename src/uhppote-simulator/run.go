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

type handler struct {
	msgType    byte
	factory    func() interface{}
	dispatcher func(*net.UDPConn, *net.UDPAddr, interface{})
}

var f = func() interface{} {
	return new(uhppote.FindDevicesRequest)
}

var handlers = map[byte]*handler{
	0x94: &handler{0x94, func() interface{} { return new(uhppote.FindDevicesRequest) }, find},
	0x50: &handler{0x50, func() interface{} { return new(uhppote.PutCardRequest) }, putCard},
	0x52: &handler{0x52, func() interface{} { return new(uhppote.DeleteCardRequest) }, deleteCard},
	0x58: &handler{0x58, func() interface{} { return new(uhppote.GetCardsRequest) }, getCards},
	0x5a: &handler{0x5a, func() interface{} { return new(uhppote.GetCardByIdRequest) }, getCardById},
	0x5c: &handler{0x5c, func() interface{} { return new(uhppote.GetCardByIndexRequest) }, getCardByIndex},
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

	request := h.factory()

	err := codec.Unmarshal(bytes, request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	h.dispatcher(c, src, request)
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

func find(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	for _, s := range simulators {
		response, err := s.Find(request)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			send(c, src, response)
		}
	}
}

func putCard(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	rq := request.(*uhppote.PutCardRequest)
	for _, s := range simulators {
		if s.SerialNumber == rq.SerialNumber {
			response, err := s.PutCard(rq)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func deleteCard(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	rq := request.(*uhppote.DeleteCardRequest)
	for _, s := range simulators {
		if s.SerialNumber == rq.SerialNumber {
			response, err := s.DeleteCard(rq)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func getCards(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	rq := request.(*uhppote.GetCardsRequest)
	for _, s := range simulators {
		if s.SerialNumber == rq.SerialNumber {
			response, err := s.GetCards(rq)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else {
				send(c, src, response)
			}
		}
	}
}

func getCardById(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	rq := request.(*uhppote.GetCardByIdRequest)
	for _, s := range simulators {
		response, err := s.GetCardById(rq)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			send(c, src, response)
		}
	}
}

func getCardByIndex(c *net.UDPConn, src *net.UDPAddr, request interface{}) {
	rq := request.(*uhppote.GetCardByIndexRequest)
	for _, s := range simulators {
		response, err := s.GetCardByIndex(rq)
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
