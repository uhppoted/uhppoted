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
	DateTime time.Time
}

func (d *DateTime) String() string {
	return d.DateTime.Format("2006-01-02 15:04:05")
}

func DecodeDateTime(bytes []byte) (*DateTime, error) {
	decoded, err := bcd.Decode(bytes)
	if err != nil {
		return nil, err
	}

	datetime, err := time.ParseInLocation("20060102150405", decoded, time.Local)
	if err != nil {
		return nil, err
	}

	return &DateTime{datetime}, nil
}

func (d *DateTime) Encode(bytes []byte) {
	encoded, err := bcd.Encode(d.DateTime.Format("20060102150405"))

	if err != nil {
		panic(fmt.Sprintf("Unexpected error encoding date-time %v to BCD: [%v]", d, err))
	} else {
		copy(bytes, *encoded)
	}
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
