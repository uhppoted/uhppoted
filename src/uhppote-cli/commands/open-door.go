package commands

import (
	"errors"
	"fmt"
	"uhppote"
)

type OpenDoorCommand struct {
	SerialNumber uint32
	Door         byte
}

func NewOpenDoorCommand() (*OpenDoorCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	door, err := getUint32(2, "Missing door ID", "Invalid door ID: %v")
	if err != nil {
		return nil, err
	}

	if door != 1 && door != 2 && door != 3 && door != 4 {
		return nil, errors.New(fmt.Sprintf("Invalid door ID: %v", door))
	}

	return &OpenDoorCommand{serialNumber, byte(door)}, nil
}

func (c *OpenDoorCommand) Execute(u *uhppote.UHPPOTE) error {
	opened, err := u.OpenDoor(c.SerialNumber, c.Door)

	if err == nil {
		fmt.Printf("%v\n", opened)
	}

	return err
}

func (c *OpenDoorCommand) CLI() string {
	return "open"
}

func (c *OpenDoorCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] open <serial number> <door>")
	fmt.Println()
	fmt.Println(" Opens the requested door:")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <door>           (required) door to open [1,2,3 or 4]")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli open 12345678 2")
	fmt.Println()
}
