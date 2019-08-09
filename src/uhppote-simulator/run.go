package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"regexp"
	"uhppote-simulator/rest"
	"uhppote-simulator/simulator"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
)

func simulate(ctx *simulator.Context) {
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

	go run(ctx, connection, wait)

	<-interrupt
	connection.Close()
	<-wait

	os.Exit(1)
}

func run(ctx *simulator.Context, connection *net.UDPConn, wait chan int) {
	go func() {
		err := listenAndServe(ctx, connection)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0
	}()

	go func() {
		for {
			msg := ctx.DeviceList.GetMessage()
			send(connection, msg.Destination, msg.Message)
		}
	}()

	go func() {
		rest.Run(ctx)
	}()

}

func listenAndServe(ctx *simulator.Context, c *net.UDPConn) error {
	for {
		request, remote, err := receive(c)
		if err != nil {
			return err
		}

		handle(ctx, c, remote, request)
	}

	return nil
}

func handle(ctx *simulator.Context, c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	request, err := messages.UnmarshalRequest(bytes)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	f := func(s *simulator.Simulator) {
		s.Handle(src, request)
	}

	ctx.DeviceList.Apply(f)
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
