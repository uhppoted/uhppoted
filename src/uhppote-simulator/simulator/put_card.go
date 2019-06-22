package simulator

import (
	"uhppote"
	"uhppote-simulator/simulator/entities"
)

func (s *Simulator) PutCard(request uhppote.PutCardRequest) (*uhppote.PutCardResponse, error) {
	card := entities.Card{
		CardNumber: request.CardNumber,
		From:       request.From,
		To:         request.To,
		Door1:      request.Door1,
		Door2:      request.Door2,
		Door3:      request.Door3,
		Door4:      request.Door4,
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
