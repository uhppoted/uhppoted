package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"uhppote"
	"uhppote-cli/commands"
)

type bind struct {
	address net.UDPAddr
}

var cli = []commands.Command{
	&commands.HelpCommand{},
	&commands.VersionCommand{VERSION},
	&commands.GetDevicesCommand{},
	&commands.GetDoorDelayCommand{},
	&commands.SetDoorDelayCommand{},
	&commands.GetEventsCommand{},
	&commands.GetEventIndexCommand{},
	&commands.SetEventIndexCommand{},
}

var VERSION = "v0.00.0"
var debug = false
var local = bind{net.UDPAddr{net.IPv4(0, 0, 0, 0), 60001, ""}}

func main() {
	flag.Var(&local, "bind", "Sets the local IP address and port to which to bind (e.g. 192.168.0.100:60001)")
	flag.BoolVar(&debug, "debug", false, "Displays vaguely useful information while processing a command")
	flag.Parse()

	u := uhppote.UHPPOTE{
		BindAddress: local.address,
		Debug:       debug,
	}

	command, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	err = command.Execute(&u)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}

func parse() (commands.Command, error) {
	var cmd commands.Command = commands.NewHelpCommand()
	var err error = nil

	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		case "get-status":
			cmd, err = commands.NewGetStatusCommand()

		case "get-time":
			cmd, err = commands.NewGetTimeCommand()

		case "set-time":
			cmd, err = commands.NewSetTimeCommand()

		case "set-ip-address":
			cmd, err = commands.NewSetAddressCommand()

		case "get-cards":
			cmd, err = commands.NewGetCardsCommand()

		case "get-card":
			cmd, err = commands.NewGetCardCommand()

		case "grant":
			cmd, err = commands.NewGrantCommand()

		case "revoke":
			cmd, err = commands.NewRevokeCommand()

		case "revoke-all":
			cmd, err = commands.NewRevokeAllCommand()

		case "open":
			cmd, err = commands.NewOpenDoorCommand()

		default:
			for _, c := range cli {
				if c.CLI() == flag.Arg(0) {
					cmd = c
				}
			}
		}
	}

	return cmd, err
}

func (b *bind) String() string {
	return b.address.String()
}

func (b *bind) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = *address

	return nil
}
