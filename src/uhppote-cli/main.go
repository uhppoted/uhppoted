package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uhppote"
	"uhppote-cli/commands"
)

var VERSION = "v0.00.0"
var debug = false

func main() {
	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	command := parse()

	err := command.Execute()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func parse() commands.Command {
	var cmd commands.Command = commands.NewHelpCommand(debug)
	var err error = nil

	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		case "help":
			cmd = commands.NewHelpCommand(debug)

		case "version":
			cmd = commands.NewVersionCommand(VERSION, debug)

		case "list-devices":
			cmd, err = commands.NewListDevicesCommand(debug)

		case "get-authorised":
			cmd, err = commands.NewGetAuthorisedCommand(debug)

		case "get-swipe":
			cmd, err = commands.NewGetSwipeCommand(debug)
		}
	}

	if err == nil {
		return cmd
	}

	return commands.NewHelpCommand(debug)
}

func parsex(s string) func() error {
	switch s {
	case "get-time":
		return gettime

	case "set-time":
		return settime

	case "set-ip-address":
		return setaddress

	case "authorise":
		return authorise
	}

	return nil
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

	u := uhppote.UHPPOTE{}
	u.Debug = debug
	datetime, err := uhppote.GetTime(uint32(serialNumber), &u)

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
		if flag.Arg(2) == "now" {
			datetime = time.Now()
		} else {
			datetime, err = time.Parse("2006-01-02 15:04:05", flag.Arg(2))
			if err != nil {
				return errors.New(fmt.Sprintf("Invalid date/time parameter: %v", flag.Arg(3)))
			}
		}
	}

	u := uhppote.UHPPOTE{}
	u.Debug = debug
	devicetime, err := uhppote.SetTime(uint32(serialNumber), datetime, &u)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}

func setaddress() error {
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

	if len(flag.Args()) < 3 {
		return errors.New("Missing IP address")
	}

	address := net.ParseIP(flag.Arg(2))

	if address == nil || address.To4() == nil {
		return errors.New(fmt.Sprintf("Invalid IP address: %v", flag.Arg(2)))
	}

	mask := net.IPv4(255, 255, 255, 0)
	if len(flag.Args()) > 3 {
		mask = net.ParseIP(flag.Arg(3))

		if mask == nil || mask.To4() == nil {
			mask = net.IPv4(255, 255, 255, 0)
		}
	}

	gateway := net.IPv4(0, 0, 0, 0)

	if len(flag.Args()) > 4 {
		gateway = net.ParseIP(flag.Arg(3))
		if gateway == nil || gateway.To4() == nil {
			gateway = net.IPv4(0, 0, 0, 0)
		}
	}

	address = net.IPv4(192, 168, 1, 125)
	mask = net.IPv4(255, 255, 255, 0)
	gateway = net.IPv4(0, 0, 0, 0)

	u := uhppote.UHPPOTE{}
	u.Debug = debug
	err = uhppote.SetAddress(uint32(serialNumber), address, mask, gateway, &u)

	return err
}

func authorise() error {
	serialNumber, err := getSerialNumber()
	if err != nil {
		return err
	}

	cardNumber, err := getCardNumber(2)
	if err != nil {
		return err
	}

	from, err := getDate(3, "Missing start date", "Invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(4, "Missing end date", "Invalid end date: %v")
	if err != nil {
		return err
	}

	permissions, err := getPermissions(5)
	if err != nil {
		return err
	}

	u := uhppote.UHPPOTE{}
	u.SerialNumber = serialNumber
	u.Debug = debug
	authorised, err := u.Authorise(cardNumber, *from, *to, *permissions)

	if err == nil {
		fmt.Printf("%v\n", authorised)
	}

	return err
}

func getSerialNumber() (uint32, error) {
	if len(flag.Args()) < 2 {
		return 0, errors.New("Missing serial number")
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(1))

	if !valid {
		return 0, errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	serialNumber, err := strconv.ParseUint(flag.Arg(1), 10, 32)

	if err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid serial number: %v", flag.Arg(1)))
	}

	return uint32(serialNumber), err
}

func getCardNumber(index int) (uint32, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New("Missing card number")
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf("Invalid card number: %v", flag.Arg(index)))
	}

	cardNumber, err := strconv.ParseUint(flag.Arg(index), 10, 32)

	if err != nil {
		return 0, errors.New(fmt.Sprintf("Invalid card number: %v", flag.Arg(index)))
	}

	return uint32(cardNumber), err
}

func getDate(index int, missing, invalid string) (*time.Time, error) {
	if len(flag.Args()) < index+1 {
		return nil, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", flag.Arg(index))

	if !valid {
		return nil, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	date, err := time.Parse("2006-01-02", flag.Arg(index))

	if err != nil {
		return nil, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return &date, err
}

func getPermissions(index int) (*[]int, error) {
	permissions := []int{}

	if len(flag.Args()) > index {
		matches := strings.Split(flag.Arg(index), ",")

		for _, match := range matches {
			door, err := strconv.Atoi(match)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid door '%v'", match))
			}

			permissions = append(permissions, door)
		}
	}

	return &permissions, nil
}

func getUint32(index int, missing, invalid string) (uint32, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 32)

	if err != nil {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return uint32(N), err
}
