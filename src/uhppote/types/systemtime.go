package types

import (
	"time"
	"uhppote/encoding/bcd"
)

type SystemTime time.Time

func (t SystemTime) String() string {
	return time.Time(t).Format("15:04:05")
}

func (t *SystemTime) UnmarshalUT0311L0x(bytes []byte) error {
	decoded, err := bcd.Decode(bytes[0:3])
	if err != nil {
		return err
	}

	time, err := time.ParseInLocation("150405", decoded, time.Local)
	if err != nil {
		return err
	}

	*t = SystemTime(time)

	return nil
}
