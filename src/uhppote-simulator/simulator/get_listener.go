package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetListener(request *messages.GetListenerRequest) (*messages.GetListenerResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	response := messages.GetListenerResponse{
		SerialNumber: s.SerialNumber,
		Address:      s.Listener.IP,
		Port:         uint16(s.Listener.Port),
	}

	return &response, nil
}
