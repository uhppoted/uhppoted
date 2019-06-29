package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uhppote/encoding/bcd"
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) MarshalUT0311L0x() ([]byte, error) {
	encoded, err := bcd.Encode(time.Time(d).Format("20060102"))

	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Error encoding date %v to BCD: [%v]", d, err))
	}

	if encoded == nil {
		return []byte{}, errors.New(fmt.Sprintf("Unknown error encoding date %v to BCD", d))
	}

	return *encoded, nil
}

func (d *Date) UnmarshalUT0311L0x(bytes []byte) error {
	decoded, err := bcd.Decode(bytes[0:4])
	if err != nil {
		return err
	}

	date, err := time.ParseInLocation("20060102", decoded, time.Local)
	if err != nil {
		return err
	}

	*d = Date(date)

	return nil
}

func (d *Date) UnmarshalPtr(bytes []byte) (interface{}, error) {
	decoded, err := bcd.Decode(bytes[0:4])
	if err != nil {
		return nil, err
	}

	date, err := time.ParseInLocation("20060102", decoded, time.Local)
	if err != nil {
		return nil, err
	}

	x := Date(date)
	return &x, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	date, err := time.ParseInLocation("2006-01-02", s, time.Local)
	if err != nil {
		return err
	}

	*d = Date(date)

	return nil
}
