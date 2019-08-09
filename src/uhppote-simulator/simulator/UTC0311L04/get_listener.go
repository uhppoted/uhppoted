package UTC0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UTC0311L04) getListener(addr *net.UDPAddr, request *messages.GetListenerRequest) {
	if s.SerialNumber == request.SerialNumber {
		response := messages.GetListenerResponse{
			SerialNumber: s.SerialNumber,
			Address:      s.Listener.IP,
			Port:         uint16(s.Listener.Port),
		}

		s.send(addr, &response)
	}
}
