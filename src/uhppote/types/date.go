package types

import (
	"encoding/bcd"
	"encoding/json"
	"fmt"
	"time"
)

type Date time.Time

func (d Date) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (d Date) Encode(bytes []byte) {
	encoded, err := bcd.Encode(time.Time(d).Format("20060102"))

	if err != nil {
		panic(fmt.Sprintf("Unexpected error encoding date %v to BCD: [%v]", d, err))
	} else {
		copy(bytes, *encoded)
	}
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
