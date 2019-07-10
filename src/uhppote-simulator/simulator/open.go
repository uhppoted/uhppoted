package simulator

import (
	"time"
	"uhppote"
	"uhppote-simulator/simulator/entities"
	"uhppote/types"
)

func (s *Simulator) OpenDoor(request *uhppote.OpenDoorRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	granted := false
	saved := false

	if request.Door > 0 && request.Door <= 4 {
		granted = true
		utc := time.Now().UTC()
		datetime := utc.Add(time.Duration(s.TimeOffset))

		event := entities.Event{
			Type:       0x02,
			Granted:    true,
			Door:       request.Door,
			DoorOpened: true,
			UserId:     3922570474,
			Timestamp:  types.DateTime(datetime),
			RecordType: 0x2c,
		}

		s.Events.Add(&event)

		err := s.Save()
		if err == nil {
			saved = true
		}

	}

	response := uhppote.OpenDoorResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    granted && saved,
	}

	return &response, nil
}
