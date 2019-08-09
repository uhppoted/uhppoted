package UTC0311L04

import (
	"net"
	"uhppote-simulator/entities"
	"uhppote/messages"
)

func (s *UTC0311L04) putCard(addr *net.UDPAddr, request *messages.PutCardRequest) {
	if request.SerialNumber == s.SerialNumber {
		card := entities.Card{
			CardNumber: request.CardNumber,
			From:       request.From,
			To:         request.To,
			Doors: map[uint8]bool{1: request.Door1,
				2: request.Door2,
				3: request.Door3,
				4: request.Door4,
			},
		}

		s.Cards.Put(&card)

		response := messages.PutCardResponse{
			SerialNumber: s.SerialNumber,
			Succeeded:    true,
		}

		s.send(addr, &response)

		err := s.Save()
		if err == nil {
			s.onError(err)
		}
	}
}
