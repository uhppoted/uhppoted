package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) SetEventIndex(request *messages.SetEventIndexRequest) (*messages.SetEventIndexResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	if request.MagicWord != 0x55aaaa55 {
		return nil, nil
	}

	set := false
	saved := false

	if err := s.Events.SetIndex(request.Index); err == nil {
		set = true
		if err := s.Save(); err == nil {
			saved = true
		}
	}

	response := messages.SetEventIndexResponse{
		SerialNumber: s.SerialNumber,
		Success:      set && saved,
	}

	return &response, nil
}
