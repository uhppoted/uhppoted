package uhppote

import (
	"errors"
	"fmt"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) PutCard(serialNumber uint32, card types.Card) (*types.Result, error) {
	request := messages.PutCardRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   card.CardNumber,
		From:         card.From,
		To:           card.To,
		Door1:        card.Doors[0],
		Door2:        card.Doors[1],
		Door3:        card.Doors[2],
		Door4:        card.Doors[3],
	}

	reply := messages.PutCardResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	if uint32(reply.SerialNumber) != serialNumber {
		return nil, errors.New(fmt.Sprintf("Incorrect serial number in response - expect '%v', received '%v'", serialNumber, reply.SerialNumber))
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
