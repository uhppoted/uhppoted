package uhppote

import (
	"uhppote/types"
)

type DeleteCardsRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	MagicNumber  uint32             `uhppote:"offset:8"`
}

type DeleteCardsResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) DeleteCards(serialNumber uint32) (*types.Result, error) {
	request := DeleteCardsRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		MagicNumber:  0x55aaaa55,
	}

	reply := DeleteCardsResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
