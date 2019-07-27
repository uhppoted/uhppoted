package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetCardById(request *messages.GetCardByIdRequest) (*messages.GetCardByIdResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	response := messages.GetCardByIdResponse{
		SerialNumber: s.SerialNumber,
	}

	for _, card := range s.Cards {
		if request.CardNumber == card.CardNumber {
			response.CardNumber = card.CardNumber
			response.From = &card.From
			response.To = &card.To
			response.Door1 = card.Door1
			response.Door2 = card.Door2
			response.Door3 = card.Door3
			response.Door4 = card.Door4
			break
		}
	}

	return &response, nil
}

func (s *Simulator) GetCardByIndex(request *messages.GetCardByIndexRequest) (*messages.GetCardByIndexResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	response := messages.GetCardByIndexResponse{
		SerialNumber: s.SerialNumber,
	}

	if request.Index > 0 && request.Index <= uint32(len(s.Cards)) {
		card := s.Cards[request.Index-1]
		response.CardNumber = card.CardNumber
		response.From = &card.From
		response.To = &card.To
		response.Door1 = card.Door1
		response.Door2 = card.Door2
		response.Door3 = card.Door3
		response.Door4 = card.Door4
	}

	return &response, nil
}
