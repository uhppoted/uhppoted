package types

import "fmt"

type DoorDelay struct {
	SerialNumber SerialNumber
	Door         uint8
	Delay        uint8
}

type Opened struct {
	SerialNumber SerialNumber
	Door         uint32
	Opened       bool
}

func (d *DoorDelay) String() string {
	return fmt.Sprintf("%s %v %v", d.SerialNumber, d.Door, d.Delay)
}

func (r *Opened) String() string {
	return fmt.Sprintf("%s %v %v", r.SerialNumber, r.Door, r.Opened)
}
