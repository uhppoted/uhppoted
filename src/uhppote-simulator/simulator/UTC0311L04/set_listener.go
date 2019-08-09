package UTC0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UTC0311L04) setListener(addr *net.UDPAddr, request *messages.SetListenerRequest) {
	if s.SerialNumber == request.SerialNumber {

		listener := net.UDPAddr{IP: request.Address, Port: int(request.Port)}
		s.Listener = &listener

		response := messages.SetListenerResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    true,
		}

		s.send(addr, &response)

		err := s.Save()
		if err != nil {
			s.onError(err)
		}
	}
}
