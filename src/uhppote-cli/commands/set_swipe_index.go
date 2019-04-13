package commands

import (
	"fmt"
	"uhppote"
)

type SetSwipeIndexCommand struct {
}

func (c *SetSwipeIndexCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	index, err := getUint32(2, "Missing swipe index", "Invalid swipe index: %v")
	if err != nil {
		return err
	}

	result, err := u.SetSwipeIndex(serialNumber, index)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", result)

	return nil
}

func (c *SetSwipeIndexCommand) CLI() string {
	return "set-swipe-index"
}

func (c *SetSwipeIndexCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-swipe-index <serial number> <index>")
	fmt.Println()
	fmt.Println(" Sets the swipe index")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  index          (required) swipe index")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-swipe-index 12345678 15")
	fmt.Println()
}
