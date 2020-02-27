package uhppote

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

func (u *UHPPOTE) DeleteCard(serialNumber, cardNumber uint32) (*types.Result, error) {
	request := messages.DeleteCardRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   cardNumber,
	}

	reply := messages.DeleteCardResponse{}

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
