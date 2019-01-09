package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
	"uhppote"
)

var debug = false
var VERSION = "v0.00.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	cmd := flag.Arg(0)
	f := parse(cmd)

	if f == nil {
		usage()
		return
	}

	err := f()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func usage() error {
	fmt.Println("Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help     Displays this message")
	fmt.Println("             For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println("    version  Displays the current version")
	fmt.Println("    search   Searches for UHPPOTE controllers on the network")
	fmt.Println("    get-time Returns the current time on the selected controller")
	fmt.Println("    set-time Sets the current time on the selected controller")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()

	return nil
}

func parse(s string) func() error {
	switch s {
	case "help":
		return help

	case "version":
		return version

	case "search":
		return search

	case "get-time":
		return gettime

	case "set-time":
		return settime
	}

	return nil
}

func help() error {
	if len(flag.Args()) > 1 {
		switch flag.Arg(1) {
		case "commands":
			helpCommands()

		case "version":
			helpVersion()

		case "search":
			helpSearch()

		case "get-time":
			helpGetTime()

		case "set-time":
			helpSetTime()

		default:
			return errors.New(fmt.Sprintf("Invalid command: %v. Type 'help commands' to get a list of supported commands", flag.Arg(1)))
		}

	}

	return nil
}

func version() error {
	fmt.Printf("%v\n", VERSION)

	return nil
}

func search() error {
	devices, err := uhppote.Search(debug)

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}

func gettime() error {
	if len(flag.Args()) < 2 {
		return errors.New("Missing serial number")
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(1))

	if !valid {
		return errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	serialNumber, err := strconv.ParseUint(flag.Arg(1), 10, 32)

	if err != nil {
		return errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	datetime, err := uhppote.GetTime(uint32(serialNumber), debug)

	if err == nil {
		fmt.Printf("%s\n", datetime)
	}

	return err
}

func settime() error {
	if len(flag.Args()) < 2 {
		return errors.New("Missing serial number")
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(1))

	if !valid {
		return errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	serialNumber, err := strconv.ParseUint(flag.Arg(1), 10, 32)

	if err != nil {
		return errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	datetime := time.Now()

	if len(flag.Args()) > 2 {
		fmt.Printf("%v\n", flag.Args())
		switch flag.Arg(2) {
		case "-now":
			datetime = time.Now()
		case "-time":
			if len(flag.Args()) < 4 {
				return errors.New(fmt.Sprintf("Missing date/time parameter"))
			}
			datetime, err = time.Parse("2006-01-02 15:04:05", flag.Arg(3))
			if err != nil {
				return errors.New(fmt.Sprintf("Invalid date/time parameter: %v", flag.Arg(3)))
			}
		}
	}

	devicetime, err := uhppote.SetTime(uint32(serialNumber), datetime, debug)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}
