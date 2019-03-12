package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"uhppote-simulator/simulator"
	"uhppote/messages"
)

var debug = false
var VERSION = "v0.00.0"

func main() {
	flag.BoolVar(&debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	s := simulator.Simulator{Debug: debug}
	err := run(&s)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func usage() error {
	fmt.Println("Usage: uhppote-simulator [options]")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --help     Displays this message")
	fmt.Println("               For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println("    --version  Displays the current version of uhppote-simulator")
	fmt.Println("    --debug    Displays a trace of simulator activity")
	fmt.Println()

	return nil
}

func run(s *simulator.Simulator) error {
	request := make([]byte, 2048)
	local, err := net.ResolveUDPAddr("udp", "0.0.0.0:60000")

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to resolve UDP local address [%v]", err))
	}

	connection, err := net.ListenUDP("udp", local)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open UDP socket [%v]", err))
	}

	defer close(connection)

	for {
		N, remote, err := connection.ReadFromUDP(request)

		if err != nil {
			return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
		}

		if s.Debug {
			regex := regexp.MustCompile("(?m)^(.*)")

			fmt.Printf(" ... received %v bytes from %v\n", N, remote)
			fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(request[:N]), " ... $1"))
		}

		// ...parse request

		reply, err := s.Handle(request[:N])

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else {
			N, err = connection.WriteTo(reply, remote)

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

func close(connection net.Conn) {
	fmt.Println(" ... closing connection")

	connection.Close()
}

func version() error {
	fmt.Printf("%v\n", VERSION)

	return nil
}

func parse(bytes []byte) messages.Message {
	fmt.Printf("%v %x %x\n", len(bytes), bytes[0], bytes[1])
	if len(bytes) == 64 && bytes[0] == 0x17 {
		switch bytes[1] {
		case 0x94:
			return &messages.Search{}
		}
	}

	return nil
}
