package simulator

import (
	"net"
	"uhppote/messages"
)

func (s *Simulator) getListener(addr *net.UDPAddr, request *messages.GetListenerRequest) {
	if s.SerialNumber == request.SerialNumber {
		response := messages.GetListenerResponse{
			SerialNumber: s.SerialNumber,
			Address:      s.Listener.IP,
			Port:         uint16(s.Listener.Port),
		}

		s.send(addr, &response)
	}
}
