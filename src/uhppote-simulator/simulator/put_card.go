package simulator

import (
	"uhppote"
	"uhppote-simulator/simulator/entities"
)

func (s *Simulator) PutCard(request *uhppote.PutCardRequest) (*uhppote.PutCardResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
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

	response := uhppote.PutCardResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    saved,
	}

	return &response, nil
}
