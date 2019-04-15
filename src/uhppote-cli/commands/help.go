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
	&GetDevicesCommand{},
	&GetStatusCommand{},
	&GetTimeCommand{},
	&GetCardsCommand{},
	&GetCardCommand{},
	&GetEventsCommand{},
	&GrantCommand{},
	&RevokeCommand{},
	&RevokeAllCommand{},
	&GetDoorDelayCommand{},
	&SetAddressCommand{},
	&SetTimeCommand{},
	&SetDoorDelayCommand{},
	&OpenDoorCommand{},
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

func (c *HelpCommand) Execute(u *uhppote.UHPPOTE) error {
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {

			if flag.Arg(1) == "commands" {
				helpCommands()
				return nil
			}

			for _, c := range commands {
				if c.CLI() == flag.Arg(1) {
					c.Help()
					return nil
				}
			}
		}

		return errors.New(fmt.Sprintf("Invalid command: %v. Type 'help commands' to get a list of supported commands", flag.Arg(1)))
	}

	(&HelpCommand{}).Help()

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
	fmt.Println(" get-address    <serial number>")
	fmt.Println(" get-time       <serial number>")
	fmt.Println(" get-swipes     <serial number> <index>")
	fmt.Println(" set-time       <serial number> <date> <time>")
	fmt.Println(" set-ip-address <serial number> <address>")
	fmt.Println(" get-authorised <serial number>")
	fmt.Println()
}
