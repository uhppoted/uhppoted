package simulator

import (
	"net"
	"uhppote/messages"
)

func (s *Simulator) SetListener(request *messages.SetListenerRequest) (*messages.SetListenerResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	listener := net.UDPAddr{IP: request.Address, Port: int(request.Port)}
	s.Listener = &listener

	saved := false
	err := s.Save()
	if err == nil {
		saved = false
	}

	response := messages.SetListenerResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    saved,
	}

	return &response, nil
}
