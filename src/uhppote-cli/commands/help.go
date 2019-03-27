package commands

import (
	"errors"
	"flag"
	"fmt"
	"uhppote"
)

type HelpCommand struct {
}

var commands = []Command{
	&VersionCommand{},
	&OpenDoorCommand{},
	&GrantCommand{},
	&RevokeCommand{},
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (c *HelpCommand) Execute(u *uhppote.UHPPOTE) error {
	var cmd Command = nil

	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {
			switch flag.Arg(1) {
			case "commands":
				helpCommands()

			case "list-devices":
				cmd = &ListDevicesCommand{}

			case "get-status":
				cmd = &GetStatusCommand{}

			case "get-time":
				cmd = &GetTimeCommand{}

			case "set-time":
				cmd = &SetTimeCommand{}

			case "set-address":
				cmd = &SetAddressCommand{}

			case "get-authorised":
				cmd = &GetAuthorisedCommand{}

			case "get-swipe":
				cmd = &GetSwipesCommand{}

			case "authorise":

			default:
				for _, c := range commands {
					if c.CLI() == flag.Arg(1) {
						cmd = c
					}
				}
			}
		}

		if cmd == nil {
			return errors.New(fmt.Sprintf("Invalid command: %v. Type 'help commands' to get a list of supported commands", flag.Arg(1)))
		}
	}

	if cmd == nil {
		cmd = &HelpCommand{}
	}

	cmd.Help()

	return nil
}

func (c *HelpCommand) CLI() string {
	return "help"
}

func (c *HelpCommand) Help() {
	usage()
}

func usage() error {
	fmt.Println()
	fmt.Println("  Usage: uhppote-cli [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help            Displays this message")
	fmt.Println("                    For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println("    version         Displays the current version")
	fmt.Println("    list-devices    Returns a list of found UHPPOTE controllers on the network")
	fmt.Println("    get-status      Returns the current status for the selected controller")
	fmt.Println("    get-time        Returns the current time on the selected controller")
	fmt.Println("    get-swipe       Retrieves a card swipe")
	fmt.Println("    get-authorised  Retrieves list of authorised cards")
	fmt.Println("    set-ip-address  Sets the IP address on the selected controller")
	fmt.Println("    set-time        Sets the current time on the selected controller")
	fmt.Println("    set-authorised  Adds a card to the authorised cards list")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -bind   Sets the local IP address+port to use")
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()

	return nil
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()
	fmt.Println(" version")
	fmt.Println(" list-devices")
	fmt.Println(" get-status     <serial number>")
	fmt.Println(" get-time       <serial number>")
	fmt.Println(" get-swipe      <serial number> <index>")
	fmt.Println(" set-time       <serial number> <date> <time>")
	fmt.Println(" set-ip-address <serial number> <address>")
	fmt.Println(" get-authorised <serial number>")
	fmt.Println()
}

func helpListSwipes() {
	fmt.Println("Usage: uhppote-cli [options] list-swipes <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the list of recorded card swipes")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli list-swipesc 12345678")
	fmt.Println()
}

func helpAuthorise() {
	fmt.Println("Usage: uhppote-cli [options] authorise <serial number> <card number> <start date> <end date> <doors>")
	fmt.Println()
	fmt.Println(" Adds a card to the authorised list")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <doors>          (required) list of permitted doors [1 2 3 4]")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli authorise 12345678 918273645 2019-01-01 2019-12-31 1,2,4")
	fmt.Println()
}
