package uhppote

import (
	"errors"
	"fmt"
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) PutCard(serialNumber, cardNumber uint32, from, to time.Time, door1, door2, door3, door4 bool) (*types.Result, error) {
	request := messages.PutCardRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   cardNumber,
		From:         types.Date(from),
		To:           types.Date(to),
		Door1:        door1,
		Door2:        door2,
		Door3:        door3,
		Door4:        door4,
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
