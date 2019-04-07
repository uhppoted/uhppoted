package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

var debug = false
var VERSION = "v0.00.0"

func main() {
	flag.BoolVar(&debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	mac, _ := net.ParseMAC("00:66:19:39:55:2d")
	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)

	s := simulator.Simulator{
		Interrupt:    make(chan int, 1),
		Debug:        debug,
		SerialNumber: 1234567890,
		IpAddress:    net.IPv4(192, 168, 0, 25),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(0, 0, 0, 0),
		MacAddress:   mac,
		Version:      0x0892,
		Date:         types.Date{date},
	}

	go s.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Stop()

	os.Exit(1)
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

func version() error {
	fmt.Printf("%v\n", VERSION)

	return nil
}
