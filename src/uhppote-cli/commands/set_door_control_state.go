package commands

import (
	"fmt"
)

type SetDoorControlStateCommand struct {
}

func (c *SetDoorControlStateCommand) Execute(ctx Context) error {
	states := map[string]uint8{
		"normally open":   1,
		"normally closed": 2,
		"controlled":      3,
	}

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

	control, err := getString(3, "Missing control state", "Invalid control state: %v")
	if err != nil {
		return err
	}
	if _, ok := states[control]; !ok {
		return fmt.Errorf("Invalid door control state: %s (expected 'normally open', 'normally closed' or 'controlled'", control)
	}

	state, err := ctx.uhppote.GetDoorControlState(serialNumber, door)
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.SetDoorControlState(serialNumber, door, states[control], state.Delay)
	if err != nil {
		return err
	}

	fmt.Printf("%s %v %v (%s)\n", record.SerialNumber, record.Door, record.ControlState, lookup[record.ControlState])

	return nil
}

func (c *SetDoorControlStateCommand) CLI() string {
	return "set-door-control-state"
}

func (c *SetDoorControlStateCommand) Description() string {
	return "Sets the control state (normally open, normally close or controlled) for a door"
}

func (c *SetDoorControlStateCommand) Usage() string {
	return "<serial number> <door> <state>"
}

func (c *SetDoorControlStateCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-door-control-state <serial number> <door> <state>")
	fmt.Println()
	fmt.Println(" Sets the door control state")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println("  state          (required) 'normally open','normally closed', 'controlled'")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-door-control-state 12345678 3 'normally open'")
	fmt.Println()
}
