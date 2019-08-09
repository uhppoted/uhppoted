package UTC0311L04

import (
	"time"
	"uhppote-simulator/entities"
	"uhppote/types"
)

func (s *UTC0311L04) Swipe(deviceId uint32, cardNumber uint32, door uint8) (bool, uint32) {
	granted := false
	opened := false
	eventType := uint8(0x01)
	recordType := uint8(0x06)

	if s.SerialNumber == types.SerialNumber(deviceId) {
		for _, c := range s.Cards {
			if c.CardNumber == cardNumber {
				if c.Doors[door] {
					granted = true
					opened = s.Doors[door].Open()
					eventType = 0x02
					recordType = 0x2c
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
		UserId:     cardNumber,
		Timestamp:  types.DateTime(datetime),
		RecordType: recordType,
	}

	eventId := s.add(&event)

	return granted && opened, eventId
}
