package UT0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UT0311L04) getCards(addr *net.UDPAddr, request *messages.GetCardsRequest) {
	if s.SerialNumber == request.SerialNumber {
		response := messages.GetCardsResponse{
			SerialNumber: s.SerialNumber,
			Records:      uint32(len(s.Cards)),
		}

		s.send(addr, &response)
	}
}
