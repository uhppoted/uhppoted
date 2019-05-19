package uhppote

import (
	"time"
	"uhppote/types"
)

type PutCardRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x50"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         types.Date         `uhppote:"offset:12"`
	To           types.Date         `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

type PutCardResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x50"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) PutCard(serialNumber, cardNumber uint32, from, to time.Time, door1, door2, door3, door4 bool) (*types.Result, error) {
	request := PutCardRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   cardNumber,
		From:         types.Date(from),
		To:           types.Date(to),
		Door1:        door1,
		Door2:        door2,
		Door3:        door3,
		Door4:        door4,
	}

	reply := PutCardResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
