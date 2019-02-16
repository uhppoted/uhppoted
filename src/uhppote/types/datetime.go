package types

import (
	"encoding/bcd"
	"fmt"
	"time"
)

type Date struct {
	Date time.Time
}

type DateTime struct {
	SerialNumber uint32
	DateTime     time.Time
}

func (d *DateTime) String() string {
	return fmt.Sprintf("%v %v", d.SerialNumber, d.DateTime.Format("2006-01-02 15:04:05"))
}

func (d *Date) String() string {
	return d.Date.Format("2006-01-02")
}

func (d Date) Encode(bytes []byte) {
	encoded, err := bcd.Encode(d.Date.Format("20060102"))

	if err != nil {
		panic(fmt.Sprintf("Unexpected error encoding date %v to BCD: [%v]", d, err))
	} else {
		copy(bytes, *encoded)
	}
}
