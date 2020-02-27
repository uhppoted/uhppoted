package UT0311L04

import (
	"github.com/uhppoted/uhppoted/src/uhppote-simulator/entities"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"time"
)

func (s *UT0311L04) Swipe(deviceID uint32, cardNumber uint32, door uint8) (bool, uint32) {
	granted := false
	opened := false
	eventType := uint8(0x01)
	result := uint8(0x06)

	if s.SerialNumber == types.SerialNumber(deviceID) {
		for _, c := range s.Cards {
			if c.CardNumber == cardNumber {
				if c.Doors[door] {
					granted = true
					opened = s.Doors[door].Open()
					eventType = 0x02
					result = 0x01
				}
			}
		}
	}

	datetime := time.Now().UTC().Add(time.Duration(s.TimeOffset))
	event := entities.Event{
		Type:       eventType,
		Granted:    granted,
		Door:       door,
		DoorOpened: opened,
		UserID:     cardNumber,
		Timestamp:  types.DateTime(datetime),
		Result:     result,
	}

	eventID := s.add(&event)

	return granted && opened, eventID
}
