package simulator

import (
	"uhppote"
)

func (s *Simulator) DeleteCards(request *uhppote.DeleteCardsRequest) (*uhppote.DeleteCardsResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
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

	response := uhppote.DeleteCardsResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    deleted && saved,
	}

	return &response, nil
}
