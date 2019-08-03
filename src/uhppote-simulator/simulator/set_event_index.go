package simulator

import (
	"errors"
	"fmt"
	"uhppote/messages"
)

func (s *Simulator) setEventIndex(request *messages.SetEventIndexRequest) *messages.SetEventIndexResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	if request.MagicWord != 0x55aaaa55 {
		s.onError(errors.New(fmt.Sprintf("Invalid 'magic number' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)))
		return nil
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

	return &response
}
