package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetCards(request *messages.GetCardsRequest) *messages.GetCardsResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	return &messages.GetCardsResponse{
		SerialNumber: s.SerialNumber,
		Records:      uint32(len(s.Cards)),
	}
}
