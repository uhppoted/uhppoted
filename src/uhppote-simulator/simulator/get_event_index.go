package simulator

import (
	"net"
	"uhppote/messages"
)

func (s *Simulator) getEventIndex(addr *net.UDPAddr, request *messages.GetEventIndexRequest) {
	if s.SerialNumber == request.SerialNumber {
		response := messages.GetEventIndexResponse{
			SerialNumber: s.SerialNumber,
			Index:        s.Events.Index,
		}

		s.send(addr, &response)
	}
}
