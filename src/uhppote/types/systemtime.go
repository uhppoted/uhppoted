package types

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote/encoding/bcd"
	"time"
)

type SystemTime time.Time

func (t SystemTime) String() string {
	return time.Time(t).Format("15:04:05")
}

func (d SystemTime) MarshalUT0311L0x() ([]byte, error) {
	encoded, err := bcd.Encode(time.Time(d).Format("150405"))

	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Error encoding system time %v to BCD: [%v]", d, err))
	}

	if encoded == nil {
		return []byte{}, errors.New(fmt.Sprintf("Unknown error encoding system time %v to BCD", d))
	}

	return *encoded, nil
}

func (t *SystemTime) UnmarshalUT0311L0x(bytes []byte) (interface{}, error) {
	decoded, err := bcd.Decode(bytes[0:3])
	if err != nil {
		return nil, err
	}

	time, err := time.ParseInLocation("150405", decoded, time.Local)
	if err != nil {
		return nil, err
	}

	v := SystemTime(time)

	return &v, nil
}
