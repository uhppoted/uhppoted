package simulator

import (
	"uhppote"
	"uhppote/types"
)

func (s *Simulator) GetCardById(request *uhppote.GetCardByIdRequest) (interface{}, error) {
	for _, card := range s.Cards {
		if request.CardNumber == card.CardNumber {
			response := uhppote.GetCardByIdResponse{
				SerialNumber: s.SerialNumber,
				CardNumber:   card.CardNumber,
				From:         card.From,
				To:           card.To,
				Door1:        card.Door1,
				Door2:        card.Door2,
				Door3:        card.Door3,
				Door4:        card.Door4,
			}

			return &response, nil
		}
	}

	return &struct {
		MsgType      types.MsgType      `uhppote:"value:0x5a"`
		SerialNumber types.SerialNumber `uhppote:"offset:4"`
		CardNumber   uint32             `uhppote:"offset:8"`
	}{
		SerialNumber: s.SerialNumber,
		CardNumber:   0,
	}, nil
}

func (s *Simulator) GetCardByIndex(request *uhppote.GetCardByIndexRequest) (interface{}, error) {
	if request.Index > 0 && request.Index <= uint32(len(s.Cards)) {
		card := s.Cards[request.Index-1]
		response := uhppote.GetCardByIndexResponse{
			SerialNumber: s.SerialNumber,
			CardNumber:   card.CardNumber,
			From:         card.From,
			To:           card.To,
			Door1:        card.Door1,
			Door2:        card.Door2,
			Door3:        card.Door3,
			Door4:        card.Door4,
		}

		return &response, nil
	}

	return &struct {
		MsgType      types.MsgType      `uhppote:"value:0x5c"`
		SerialNumber types.SerialNumber `uhppote:"offset:4"`
		CardNumber   uint32             `uhppote:"offset:8"`
	}{
		SerialNumber: s.SerialNumber,
		CardNumber:   0,
	}, nil
}
