package commands

import (
	"errors"
	"fmt"
)

type GetEventsCommand struct {
}

func (c *GetEventsCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	first, err := ctx.uhppote.GetEvent(serialNumber, 0)
	if err != nil {
		return err
	} else if first == nil {
		return errors.New("Failed to get 'first' event")
	}

	last, err := ctx.uhppote.GetEvent(serialNumber, 0xffffffff)
	if err != nil {
		return err
	} else if last == nil {
		return errors.New("Failed to get 'last' event")
	}

	fmt.Printf("%v  %d  %d\n", serialNumber, first.Index, last.Index)

	return nil
}

func (c *GetEventsCommand) CLI() string {
	return "get-events"
}

func (c *GetEventsCommand) Description() string {
	return "Returns the indices of the 'first' and 'last' events stored on the controller"
}

func (c *GetEventsCommand) Usage() string {
	return "<serial number>"
}

func (c *GetEventsCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-events <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the indices of the first and last' events stored in the controller event buffer")
	fmt.Println(" The controller event buffer is implemented as a ring buffer with capacity for (apparently)")
	fmt.Println(" 100000 events i.e. the index of the 'last' event may be less than the index of the 'first'")
	fmt.Println(" if the event buffer has wrapped around")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-events 12345678")
	fmt.Println()
	fmt.Println("    > 12345678  10  71")
	fmt.Println()
}
