package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
)

var debug = false
var VERSION = "v0.00.0"

func main() {
	flag.BoolVar(&debug, "--debug", false, "Displays simulator activity")
	flag.Parse()

	err := run(true)

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
	fmt.Println("    --verison  Displays the current version of uhppote-simulator")
	fmt.Println("    --debug    Displays a trace of simulator activity")
	fmt.Println()

	return nil
}

func run(debug bool) error {
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

	N, remote, err := connection.ReadFromUDP(request)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... received %v bytes from %v\n", N, remote)
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(request[:N]), " ... $1"))
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
