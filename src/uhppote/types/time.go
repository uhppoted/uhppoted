package types

import (
	"fmt"
)

type Time struct {
	SerialNumber SerialNumber
	DateTime     DateTime
}

func (t Time) String() string {
	return fmt.Sprintf("%s %s", t.SerialNumber, t.DateTime.String())
}
