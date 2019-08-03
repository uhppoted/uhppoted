package simulator

import (
	"time"
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) openDoor(request *messages.OpenDoorRequest) *messages.OpenDoorResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	granted := false
	opened := false
	door := request.Door

	if !(door < 1 || door > 4) {
		granted = true
		opened = s.Doors[door].Open()

		datetime := time.Now().UTC().Add(time.Duration(s.TimeOffset))
		event := entities.Event{
			Type:       0x02,
			Granted:    granted,
			Door:       door,
			DoorOpened: opened,
			UserId:     3922570474,
			Timestamp:  types.DateTime(datetime),
			RecordType: 0x2c,
		}

		s.add(&event)
	}

	response := messages.OpenDoorResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    granted && opened,
	}

	return &response
}
