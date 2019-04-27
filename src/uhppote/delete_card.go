package uhppote

import (
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

	err := u.Execute(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
