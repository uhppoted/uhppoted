package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetCardsRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
}

type GetCardsResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Records      uint32 `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetCards(serialNumber uint32) (*types.RecordCount, error) {
	request := GetCardsRequest{
		MsgType:      0x58,
		SerialNumber: serialNumber,
	}

	reply := GetCardsResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x58 {
		return nil, errors.New(fmt.Sprintf("GetCards returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.RecordCount{
		SerialNumber: serialNumber,
		Records:      reply.Records,
	}, nil
}
