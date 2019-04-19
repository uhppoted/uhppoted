package commands

import (
	"errors"
	"fmt"
	"uhppote"
)

type OpenDoorCommand struct {
}

func (c *OpenDoorCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	door, err := getUint32(2, "Missing door ID", "Invalid door ID: %v")
	if err != nil {
		return err
	}

	if door != 1 && door != 2 && door != 3 && door != 4 {
		return errors.New(fmt.Sprintf("Invalid door ID: %v", door))
	}

	opened, err := u.OpenDoor(serialNumber, byte(door))

	if err == nil {
		fmt.Printf("%v\n", opened)
	}

	return err
}

func (c *OpenDoorCommand) CLI() string {
	return "open"
}

func (c *OpenDoorCommand) Description() string {
	return "Opens a door"
}

func (c *OpenDoorCommand) Usage() string {
	return "<serial number> <door>"
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
