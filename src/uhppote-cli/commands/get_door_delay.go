package commands

import (
	"fmt"
)

type GetDoorDelayCommand struct {
}

func (c *GetDoorDelayCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	door, err := getDoor(2, "Missing door", "Invalid door: %v")
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.GetDoorControlState(serialNumber, door)
	if err != nil {
		return err
	}

	fmt.Printf("%s %v %v\n", record.SerialNumber, record.Door, record.Delay)

	return nil
}

func (c *GetDoorDelayCommand) CLI() string {
	return "get-door-delay"
}

func (c *GetDoorDelayCommand) Description() string {
	return "Gets the time a door lock is kept open"
}

func (c *GetDoorDelayCommand) Usage() string {
	return "<serial number> <door>"
}

func (c *GetDoorDelayCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-door-delay <serial number> <door>")
	fmt.Println()
	fmt.Println(" Retrieves the door open delay (in seconds)")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-door-delay 12345678 3")
	fmt.Println()
}
