package commands

import (
	"fmt"
	"uhppote"
)

type GetSwipeIndexCommand struct {
}

func (c *GetSwipeIndexCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := u.GetSwipeIndex(serialNumber)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", index.String())

	return nil
}

func (c *GetSwipeIndexCommand) CLI() string {
	return "get-swipe-index"
}

func (c *GetSwipeIndexCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-swipe-index <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the current swipe record index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-swipe-index 12345678")
	fmt.Println()
}
