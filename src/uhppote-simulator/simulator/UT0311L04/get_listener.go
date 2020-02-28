package UT0311L04

import (
	"github.com/uhppoted/uhppote-core/messages"
	"net"
)

func (s *UT0311L04) getListener(addr *net.UDPAddr, request *messages.GetListenerRequest) {
	address := net.IPv4(0, 0, 0, 0)
	port := uint16(0)

	if s.Listener != nil {
		address = s.Listener.IP
		port = uint16(s.Listener.Port)
	}

	if s.SerialNumber == request.SerialNumber {
		response := messages.GetListenerResponse{
			SerialNumber: s.SerialNumber,
			Address:      address,
			Port:         port,
		}

		s.send(addr, &response)
	}
}
