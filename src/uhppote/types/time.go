package types

import (
	"fmt"
)

type Time struct {
	SerialNumber uint32
	DateTime     DateTime
}

func (t Time) String() string {
	return fmt.Sprintf("%v %s", t.SerialNumber, t.DateTime.String())
}
