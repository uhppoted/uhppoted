package uhppote

import (
	"uhppote/types"
)

type GetCardByIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x5c"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

type GetCardByIdRequest struct {
	MsgType      types.MsgType `uhppote:"value:0x5a"`
	SerialNumber uint32        `uhppote:"offset:4"`
	CardNumber   uint32        `uhppote:"offset:8"`
}

type GetCardByIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x5c"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         types.Date         `uhppote:"offset:12"`
	To           types.Date         `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

type GetCardByIdResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x5a"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         types.Date         `uhppote:"offset:12"`
	To           types.Date         `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

func (u *UHPPOTE) GetCardByIndex(serialNumber, index uint32) (*types.Card, error) {
	request := GetCardByIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
	}

	reply := GetCardByIndexResponse{}

	err := u.Execute(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Card{
		SerialNumber: reply.SerialNumber,
		CardNumber:   reply.CardNumber,
		From:         reply.From,
		To:           reply.To,
		Door1:        reply.Door1,
		Door2:        reply.Door2,
		Door3:        reply.Door3,
		Door4:        reply.Door4,
	}, nil
}

func (u *UHPPOTE) GetCardById(serialNumber, cardNumber uint32) (*types.Card, error) {
	request := GetCardByIdRequest{
		SerialNumber: serialNumber,
		CardNumber:   cardNumber,
	}

	reply := GetCardByIdResponse{}

	err := u.Execute(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Card{
		SerialNumber: reply.SerialNumber,
		CardNumber:   reply.CardNumber,
		From:         reply.From,
		To:           reply.To,
		Door1:        reply.Door1,
		Door2:        reply.Door2,
		Door3:        reply.Door3,
		Door4:        reply.Door4,
	}, nil
}
