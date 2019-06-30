package types

import (
	"time"
	"uhppote/encoding/bcd"
)

type SystemTime time.Time

func (t SystemTime) String() string {
	return time.Time(t).Format("15:04:05")
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
