package simulator

import (
	"uhppote"
)

func (s *Simulator) GetCards(request *uhppote.GetCardsRequest) (*uhppote.GetCardsResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	return &uhppote.GetCardsResponse{
		SerialNumber: s.SerialNumber,
		Records:      uint32(len(s.Cards)),
	}, nil
}
