package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) deleteCards(request *messages.DeleteCardsRequest) *messages.DeleteCardsResponse {
	if request.SerialNumber != s.SerialNumber {
		return nil
	}

	deleted := false
	saved := false

	if request.MagicNumber == 0x55aaaa55 {
		if deleted = s.Cards.DeleteAll(); deleted {
			if err := s.Save(); err == nil {
				saved = true
			}
		}
	}

	response := messages.DeleteCardsResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    deleted && saved,
	}

	return &response
}
