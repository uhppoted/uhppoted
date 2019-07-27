package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetCards(request *messages.GetCardsRequest) (*messages.GetCardsResponse, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	return &messages.GetCardsResponse{
		SerialNumber: s.SerialNumber,
		Records:      uint32(len(s.Cards)),
	}, nil
}
