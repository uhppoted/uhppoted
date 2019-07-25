package simulator

import (
	"time"
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) OpenDoor(request *messages.OpenDoorRequest) (*messages.OpenDoorResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	granted := false
	door := request.Door

	if !(door < 1 || door > 4) {
		granted = s.Doors[door].Open()

		datetime := time.Now().UTC().Add(time.Duration(s.TimeOffset))
		event := entities.Event{
			Type:       0x02,
			Granted:    granted,
			Door:       door,
			DoorOpened: true,
			UserId:     3922570474,
			Timestamp:  types.DateTime(datetime),
			RecordType: 0x2c,
		}

		s.Add(&event)
	}

	response := messages.OpenDoorResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    granted,
	}

	return &response, nil
}
