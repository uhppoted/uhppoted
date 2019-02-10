package types

import (
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
	bytes[0] = bcd(d.Date.Year() / 100)
	bytes[1] = bcd(d.Date.Year() % 100)
	bytes[2] = bcd(int(d.Date.Month()))
	bytes[3] = bcd(d.Date.Day())
}

func bcd(b int) byte {
	msb := b / 10
	lsb := b % 10

	return byte(msb*16 + lsb)
}
