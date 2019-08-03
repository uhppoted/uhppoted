package simulator

import (
	"net"
	"uhppote/messages"
)

func (s *Simulator) setListener(request *messages.SetListenerRequest) *messages.SetListenerResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	listener := net.UDPAddr{IP: request.Address, Port: int(request.Port)}
	s.Listener = &listener

	saved := false
	err := s.Save()
	if err == nil {
		saved = true
	}

	response := messages.SetListenerResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    saved,
	}

	return &response
}
