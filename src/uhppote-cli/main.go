package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"uhppote"
	"uhppote-cli/commands"
)

type bind struct {
	address net.IP
	port    uint
}

var VERSION = "v0.00.0"
var debug = false
var local = bind{net.ParseIP("0.0.0.0"), 60001}

func main() {
	flag.Var(&local, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	u := uhppote.UHPPOTE{
		BindAddress: local.address,
		BindPort:    local.port,
		Debug:       debug,
	}

	command := parse(u)
	err := command.Execute()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func parse(u uhppote.UHPPOTE) commands.Command {
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

		case "get-time":
			cmd, err = commands.NewGetTimeCommand(debug)

		case "set-time":
			cmd, err = commands.NewSetTimeCommand(debug)

		case "set-ip-address":
			cmd, err = commands.NewSetAddressCommand(debug)

		case "get-authorised":
			cmd, err = commands.NewGetAuthorisedCommand(debug)

		case "get-swipe":
			cmd, err = commands.NewGetSwipeCommand(debug)

		case "authorise":
			cmd, err = commands.NewGrantCommand(debug)
		}
	}

	if err == nil {
		return cmd
	}

	return commands.NewHelpCommand(debug)
}

func (b *bind) String() string {
	return net.JoinHostPort(b.address.String(), strconv.Itoa(int(b.port)))
}

func (b *bind) Set(s string) error {
	h, p, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}

	host := net.ParseIP(h)
	if host == nil {
		return errors.New(fmt.Sprintf("Invalid bind address: %s", s))
	}

	port, err := strconv.Atoi(p)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid bind port: %s", s))
	}
	if port < 1 || port > 65535 {
		return errors.New(fmt.Sprintf("Invalid bind port: %v - expected port in range 1 to 65535", port))
	}

	b.address = host
	b.port = uint(port)

	return nil
}
