package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) DeleteCard(request *messages.DeleteCardRequest) (*messages.DeleteCardResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	deleted := s.Cards.Delete(request.CardNumber)
	saved := false

	if deleted {
		if err := s.Save(); err == nil {
			saved = true
		}
	}

	response := messages.DeleteCardResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    deleted && saved,
	}

	return &response, nil
}
