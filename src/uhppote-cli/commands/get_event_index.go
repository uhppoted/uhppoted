package commands

import (
	"fmt"
	"uhppote"
)

type GetEventIndexCommand struct {
}

func (c *GetEventIndexCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := u.GetEventIndex(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", index.String())

	return nil
}

func (c *GetEventIndexCommand) CLI() string {
	return "get-swipe-index"
}

func (c *GetEventIndexCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-event-index <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the current event record index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-event-index 12345678")
	fmt.Println()
}
