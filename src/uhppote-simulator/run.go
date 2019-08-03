package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"uhppote-simulator/rest"
	"uhppote-simulator/simulator/entities"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
)

func simulate() {
	bind, err := net.ResolveUDPAddr("udp", ":60000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP bind address [%v]", err)))
		return
	}

	connection, err := net.ListenUDP("udp", bind)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to bind to UDP socket [%v]", err)))
		return
	}

	defer connection.Close()

	wait := make(chan int, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go run(connection, wait)

	<-interrupt
	connection.Close()
	<-wait

	os.Exit(1)
}

func run(connection *net.UDPConn, wait chan int) {
	queue := make(chan entities.Message, 8)

	for _, s := range simulators {
		s.TxQueue = queue
	}

	go func() {
		err := listenAndServe(connection, queue)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0
	}()

	go func() {
		for {
			msg := <-queue
			send(connection, msg.Destination, msg.Message)
		}
	}()

	go func() {
		rest.Run(simulators)
	}()

}

func listenAndServe(c *net.UDPConn, queue chan entities.Message) error {
	for {
		request, remote, err := receive(c)
		if err != nil {
			return err
		}

		handle(c, remote, request, queue)
	}

	return nil
}

func handle(c *net.UDPConn, src *net.UDPAddr, bytes []byte, queue chan entities.Message) {
	request, err := messages.UnmarshalRequest(bytes)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response := s.Handle(bytes[1], request)
		if response != nil && !reflect.ValueOf(response).IsNil() {
			queue <- entities.Message{src, response}
		}
	}
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
	msg, err := codec.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	N, err := c.WriteTo(msg, dest)
	if err != nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err)))
	} else {
		if options.debug {
			fmt.Printf(" ... sent %v bytes to %v\n%s\n", N, dest, dump(msg[0:N], " ...          "))
		}
	}
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
