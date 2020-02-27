package types

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote/encoding/bcd"
	"time"
)

type SystemDate time.Time

func (d SystemDate) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d SystemDate) MarshalUT0311L0x() ([]byte, error) {
	encoded, err := bcd.Encode(time.Time(d).Format("060102"))

	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Error encoding system date %v to BCD: [%v]", d, err))
	}

	if encoded == nil {
		return []byte{}, errors.New(fmt.Sprintf("Unknown error encoding system date %v to BCD", d))
	}

	return *encoded, nil
}

func (d *SystemDate) UnmarshalUT0311L0x(bytes []byte) (interface{}, error) {
	decoded, err := bcd.Decode(bytes[0:3])
	if err != nil {
		return nil, err
	}

	date, err := time.ParseInLocation("060102", decoded, time.Local)
	if err != nil {
		return nil, err
	}

	v := SystemDate(date)

	return &v, nil
}
