package types

import (
	"fmt"
	"time"
)

type DateTime struct {
	SerialNumber uint32
	DateTime     time.Time
}

func (d *DateTime) String() string {
	return fmt.Sprintf("%v %v", d.SerialNumber, d.DateTime.Format("2006-01-02 15:04:05"))
}
