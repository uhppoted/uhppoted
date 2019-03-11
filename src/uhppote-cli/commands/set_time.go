package commands

import (
	"errors"
	"flag"
	"fmt"
	"time"
	"uhppote"
	"uhppote/types"
)

type SetTimeCommand struct {
	SerialNumber uint32
	DateTime     types.DateTime
}

func NewSetTimeCommand() (*SetTimeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	datetime := time.Now()
	if len(flag.Args()) > 2 {
		if flag.Arg(2) == "now" {
			datetime = time.Now()
		} else {
			datetime, err = time.Parse("2006-01-02 15:04:05", flag.Arg(2))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid date/time parameter: %v", flag.Arg(3)))
			}
		}
	}

	return &SetTimeCommand{serialNumber, types.DateTime{datetime}}, nil
}

func (c *SetTimeCommand) Execute(u *uhppote.UHPPOTE) error {
	devicetime, err := u.SetTime(c.SerialNumber, c.DateTime)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}

func (c *SetTimeCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Sets the controller date/time to the supplied time. Defaults to 'now'. Command format")
	fmt.Println()
	fmt.Println(" <serial number> [now|<yyyy-mm-dd HH:mm:ss>]")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
	fmt.Println("    now                    Sets the controller time to the system time of the local system")
	fmt.Println("    'yyyy-mm-dd HH:mm:ss'  Sets the controller time to the explicitly supplied instant")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-time")
	fmt.Println("    uhppote-cli set-time now")
	fmt.Println("    uhppote-cli set-time '2019-01-12 20:15:32'")
	fmt.Println()
}
