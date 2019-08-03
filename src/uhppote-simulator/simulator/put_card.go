package simulator

import (
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
)

func (s *Simulator) putCard(request *messages.PutCardRequest) *messages.PutCardResponse {
	if request.SerialNumber != s.SerialNumber {
		return nil
	}

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

	saved := false
	err := s.Save()
	if err == nil {
		saved = true
	}

	response := messages.PutCardResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    saved,
	}

	return &response
}
