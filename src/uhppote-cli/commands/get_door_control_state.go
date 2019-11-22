package commands

import (
	"fmt"
)

type GetDoorControlStateCommand struct {
}

func (c *GetDoorControlStateCommand) Execute(ctx Context) error {
	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

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

	fmt.Printf("%s %v %v (%s)\n", record.SerialNumber, record.Door, record.ControlState, lookup[record.ControlState])

	return nil
}

func (c *GetDoorControlStateCommand) CLI() string {
	return "get-door-control-state"
}

func (c *GetDoorControlStateCommand) Description() string {
	return "Gets the control state (normally open, normally closed or controlled) for a door"
}

func (c *GetDoorControlStateCommand) Usage() string {
	return "<serial number> <door>"
}

func (c *GetDoorControlStateCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-door-control-state <serial number> <door>")
	fmt.Println()
	fmt.Println(" Retrieves the door open delay (in seconds)")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-door-control-state 12345678 3")
	fmt.Println()
}
