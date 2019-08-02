package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetEvent(request *messages.GetEventRequest) *messages.GetEventResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	index := request.Index
	if index > s.Events.LastIndex() {
		index = s.Events.LastIndex()
	}

	if event := s.Events.Get(index); event != nil {
		response := messages.GetEventResponse{
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

		return &response
	}

	return nil
}
