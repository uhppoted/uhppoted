package UT0311L04

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/messages"
	"net"
)

func (s *UT0311L04) setListener(addr *net.UDPAddr, request *messages.SetListenerRequest) {
	if s.SerialNumber == request.SerialNumber {

		listener := net.UDPAddr{IP: request.Address, Port: int(request.Port)}
		s.Listener = &listener

		response := messages.SetListenerResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    true,
		}

		s.send(addr, &response)

		if err := s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
