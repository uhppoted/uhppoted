package types

import "fmt"

type DoorControlState struct {
	SerialNumber SerialNumber
	Door         uint8
	ControlState uint8
	Delay        uint8
}

func (d *DoorControlState) String() string {
	return fmt.Sprintf("%s %v %v %v", d.SerialNumber, d.Door, d.ControlState, d.Delay)
}
