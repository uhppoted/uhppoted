package simulator

import (
	"uhppote"
)

func (s *Simulator) GetCards(request uhppote.GetCardsRequest) (*uhppote.GetCardsResponse, error) {
	response := uhppote.GetCardsResponse{
		SerialNumber: s.SerialNumber,
		Records:      uint32(len(s.Cards)),
	}

	return &response, nil
}
