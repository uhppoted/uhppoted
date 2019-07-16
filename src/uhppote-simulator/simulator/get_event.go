package simulator

import (
	"uhppote"
)

func (s *Simulator) GetEvent(request *uhppote.GetEventRequest) (*uhppote.GetEventResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	index := request.Index
	if index > s.Events.LastIndex() {
		index = s.Events.LastIndex()
	}

	if event := s.Events.Get(index); event != nil {
		response := uhppote.GetEventResponse{
			SerialNumber: s.SerialNumber,
			Index:        index,
			Type:         event.Type,
			Granted:      event.Granted,
			Door:         event.Door,
			DoorOpened:   event.DoorOpened,
			UserId:       event.UserId,
			Timestamp:    event.Timestamp,
			RecordType:   event.RecordType,
		}

		return &response, nil
	}

	return nil, nil
}
