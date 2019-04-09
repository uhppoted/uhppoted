package types

import "fmt"

type DoorDelay struct {
	SerialNumber uint32
	Door         uint8
	Delay        uint8
}

type Opened struct {
	SerialNumber uint32
	Door         uint32
	Opened       bool
}

func (d *DoorDelay) String() string {
	return fmt.Sprintf("%v %v %v", d.SerialNumber, d.Door, d.Delay)
}

func (r *Opened) String() string {
	return fmt.Sprintf("%v %v %v", r.SerialNumber, r.Door, r.Opened)
}
