package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetCardByIndexRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Index        uint32 `uhppote:"offset:8"`
}

type GetCardByIdRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	CardNumber   uint32 `uhppote:"offset:8"`
}

type GetCardResponse struct {
	MsgType      byte       `uhppote:"offset:1"`
	SerialNumber uint32     `uhppote:"offset:4"`
	CardNumber   uint32     `uhppote:"offset:8"`
	From         types.Date `uhppote:"offset:12"`
	To           types.Date `uhppote:"offset:16"`
	Door1        bool       `uhppote:"offset:20"`
	Door2        bool       `uhppote:"offset:21"`
	Door3        bool       `uhppote:"offset:22"`
	Door4        bool       `uhppote:"offset:23"`
}

func (u *UHPPOTE) GetCardByIndex(serialNumber, index uint32) (*types.Card, error) {
	request := GetCardByIndexRequest{
		MsgType:      0x5C,
		SerialNumber: serialNumber,
		Index:        index,
	}

	reply := GetCardResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x5C {
		return nil, errors.New(fmt.Sprintf("GetCardByIndex returned incorrect message type: %02X\n", reply.MsgType))
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
		MsgType:      0x5A,
		SerialNumber: serialNumber,
		CardNumber:   cardNumber,
	}

	reply := GetCardResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x5A {
		return nil, errors.New(fmt.Sprintf("GetCardById returned incorrect message type: %02X\n", reply.MsgType))
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
