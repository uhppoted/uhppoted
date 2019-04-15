package commands

import (
	"fmt"
	"uhppote"
)

type SetEventIndexCommand struct {
}

func (c *SetEventIndexCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := getUint32(2, "Missing event index", "Invalid event index: %v")
	if err != nil {
		return err
	}

	result, err := u.SetEventIndex(serialNumber, index)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", result)

	return nil
}

func (c *SetEventIndexCommand) CLI() string {
	return "set-event-index"
}

func (c *SetEventIndexCommand) Description() string {
	return "Sets the current event index"
}

func (c *SetEventIndexCommand) Usage() string {
	return "<serial number> <index>"
}

func (c *SetEventIndexCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-event-index <serial number> <index>")
	fmt.Println()
	fmt.Println(" Sets the event index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  index          (required) event index")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-event-index 12345678 15")
	fmt.Println()
}
