package UTC0311L04

import (
	"fmt"
	"net"
	"uhppote/messages"
)

func (s *UTC0311L04) deleteCards(addr *net.UDPAddr, request *messages.DeleteCardsRequest) {
	if request.SerialNumber == s.SerialNumber {
		deleted := false

		if request.MagicWord == 0x55aaaa55 {
			deleted = s.Cards.DeleteAll()
		}

		response := messages.DeleteCardsResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    deleted,
		}

		s.send(addr, &response)

		if deleted {
			if err := s.Save(); err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}
}
