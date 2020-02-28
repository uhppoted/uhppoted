package UT0311L04

import (
	"github.com/uhppoted/uhppote-core/messages"
	"net"
)

func (s *UT0311L04) getEventIndex(addr *net.UDPAddr, request *messages.GetEventIndexRequest) {
	if s.SerialNumber == request.SerialNumber {
		response := messages.GetEventIndexResponse{
			SerialNumber: s.SerialNumber,
			Index:        s.Events.Index,
		}

		s.send(addr, &response)
	}
}
