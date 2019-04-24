package types

import "fmt"

type DoorDelay struct {
	SerialNumber SerialNumber
	Door         uint8
	Delay        uint8
}

func (d *DoorDelay) String() string {
	return fmt.Sprintf("%s %v %v", d.SerialNumber, d.Door, d.Delay)
}
