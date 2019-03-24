package types

import "fmt"

type Opened struct {
	SerialNumber uint32
	Door         uint32
	Opened       bool
}

func (r *Opened) String() string {
	return fmt.Sprintf("%v %v %v", r.SerialNumber, r.Door, r.Opened)
}
