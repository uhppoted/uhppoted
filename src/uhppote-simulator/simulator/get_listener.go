package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) getListener(request *messages.GetListenerRequest) *messages.GetListenerResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	response := messages.GetListenerResponse{
		SerialNumber: s.SerialNumber,
		Address:      s.Listener.IP,
		Port:         uint16(s.Listener.Port),
	}

	return &response
}
