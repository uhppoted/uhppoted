package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type DeleteCardRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x52"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
}

type DeleteCardResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x52"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) DeleteCard(serialNumber, cardNumber uint32) (*types.Result, error) {
	request := DeleteCardRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   cardNumber,
	}

	reply := DeleteCardResponse{}

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
