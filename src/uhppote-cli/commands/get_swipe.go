package commands

import (
	"fmt"
	"uhppote"
)

type GetSwipesCommand struct {
	SerialNumber uint32
	Index        uint32
}

func NewGetSwipesCommand() (*GetSwipesCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	index, err := getUint32(2, "Missing swipe index", "Invalid swipe index: %v")
	if err != nil {
		return nil, err
	}

	return &GetSwipesCommand{serialNumber, index}, nil
}

func (c *GetSwipesCommand) Execute(u *uhppote.UHPPOTE) error {
	swipe, err := u.GetSwipe(c.SerialNumber, c.Index)

	if err == nil {
		if swipe != nil {
			fmt.Printf("%12d %s\n", c.SerialNumber, swipe.String())
		}
	}

	return err
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
