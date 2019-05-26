package commands

import (
	"fmt"
)

type GetEventsCommand struct {
}

func (c *GetEventsCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	last, err := ctx.uhppote.GetEvent(serialNumber, 0xffffffff)
	if err != nil {
		return err
	}

	if last != nil {
		for index := last.Index; index > 0; index-- {
			swipe, err := ctx.uhppote.GetEvent(serialNumber, index)
			if err != nil {
				return err
			}

			if swipe != nil {
				fmt.Printf("%s\n", swipe.String())
			}
		}
	}

	return nil
}

func (c *GetEventsCommand) CLI() string {
	return "get-events"
}

func (c *GetEventsCommand) Description() string {
	return "Returns the list of access events stored on the controller"
}

func (c *GetEventsCommand) Usage() string {
	return "<serial number>"
}

func (c *GetEventsCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-eventsi <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the list of recorded access events")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-events 12345678")
	fmt.Println()
}
