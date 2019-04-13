package commands

import (
	"fmt"
	"uhppote"
)

type GetSwipesCommand struct {
}

func (c *GetSwipesCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	last, err := u.GetSwipe(serialNumber, 0xffffffff)
	if err != nil {
		return err
	}

	if last != nil {
		for index := last.Index; index > 0; index-- {
			swipe, err := u.GetSwipe(serialNumber, index)
			if err != nil {
				return err
			}

			if swipe != nil {
				fmt.Printf("%s\n", swipe.String())
			}
		}
	}

	//u.GetSwipeIndex(serialNumber)
	//u.SetSwipeIndex(serialNumber, 11)
	//u.GetSwipeIndex(serialNumber)

	return nil
}

func (c *GetSwipesCommand) CLI() string {
	return "get-swipes"
}

func (c *GetSwipesCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-swipes <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the list of recorded card swipes")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-swipes 12345678")
	fmt.Println()
}
