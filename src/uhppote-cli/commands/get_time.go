package commands

import (
	"fmt"
)

type GetTimeCommand struct {
}

func (c *GetTimeCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	datetime, err := ctx.uhppote.GetTime(serialNumber)

	if err == nil {
		fmt.Printf("%v\n", datetime)
	}

	return err
}

func (c *GetTimeCommand) CLI() string {
	return "get-time"
}

func (c *GetTimeCommand) Description() string {
	return "Returns the current time on the selected controller"
}

func (c *GetTimeCommand) Usage() string {
	return "<serial number>"
}

func (c *GetTimeCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Retrieves the current date/time referenced to the local timezone for the controller")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}
