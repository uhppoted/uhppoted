package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"uhppote-simulator/simulator"
)

var debug = false
var VERSION = "v0.00.0"

func main() {
	flag.BoolVar(&debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	s := simulator.Simulator{
		Interrupt: make(chan int, 1),
		Debug:     debug,
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
