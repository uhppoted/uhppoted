package types

import (
	"time"
	"uhppote/encoding/bcd"
)

type SystemDate time.Time

func (d SystemDate) String() string {
	return time.Time(d).Format("2006-01-02")
}

//func (d SystemDate) MarshalUT0311L0x() ([]byte, error) {
//	encoded, err := bcd.Encode(time.Time(d).Format("20060102"))
//
//	if err != nil {
//		return []byte{}, errors.New(fmt.Sprintf("Error encoding date %v to BCD: [%v]", d, err))
//	}
//
//	if encoded == nil {
//		return []byte{}, errors.New(fmt.Sprintf("Unknown error encoding date %v to BCD", d))
//	}
//
//	return *encoded, nil
//}

func (d *SystemDate) UnmarshalUT0311L0x(bytes []byte) error {
	decoded, err := bcd.Decode(bytes[0:3])
	if err != nil {
		return err
	}

	date, err := time.ParseInLocation("060102", decoded, time.Local)
	if err != nil {
		return err
	}

	*d = SystemDate(date)

	return nil
}
