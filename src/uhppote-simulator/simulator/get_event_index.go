package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetEventIndex(request *messages.GetEventIndexRequest) (*messages.GetEventIndexResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	response := messages.GetEventIndexResponse{
		SerialNumber: s.SerialNumber,
		Index:        s.Events.Index,
	}

	return &response, nil
}
