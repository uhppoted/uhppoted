package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"
	"uhppote/types"
)

type MacAddress net.HardwareAddr
type Version types.Version
type Date time.Time

type Card struct {
	CardNumber uint32 `json:"number"`
	From       Date   `json:"from"`
	To         Date   `json:"to"`
	Door1      bool   `json:"door-1"`
	Door2      bool   `json:"door-2"`
	Door3      bool   `json:"door-3"`
	Door4      bool   `json:"door-4"`
}

type CardList []*Card

func (m MacAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.HardwareAddr(m).String())
}

func (m *MacAddress) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	mac, err := net.ParseMAC(s)
	if err != nil {
		return err
	}

	*m = MacAddress(mac)

	return nil
}

func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%04x", v))
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	N, err := fmt.Sscanf(s, "%04x", v)
	if err != nil {
		return err
	}

	if N != 1 {
		return errors.New("Unable to extract 'version' from JSON file")
	}

	return nil
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

// TODO: implement Marshal/Unmarshal
func (l *CardList) Put(card *Card) {
	if card != nil {
		index := -1
		for i, c := range *l {
			if c.CardNumber == card.CardNumber {
				index = i
				break
			}
		}

		if index == -1 {
			*l = append(*l, card)
		} else {
			(*l)[index] = card
		}
	}
}
