package main

import (
	"flag"
	"fmt"
	"os"
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
