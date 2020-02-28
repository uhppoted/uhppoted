package UT0311L04

import (
	"github.com/uhppoted/uhppote-core/messages"
	"net"
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
