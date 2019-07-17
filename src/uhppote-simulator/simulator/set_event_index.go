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

	updated := false
	saved := false

	if updated = s.Events.SetIndex(request.Index); updated {
		if err := s.Save(); err == nil {
			saved = true
		}
	}

	response := messages.SetEventIndexResponse{
		SerialNumber: s.SerialNumber,
		Changed:      updated && saved,
	}

	return &response, nil
}
