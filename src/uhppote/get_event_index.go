package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetEventIndexRequest struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetEventIndexResponse struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetEventIndex(serialNumber uint32) (*types.EventIndex, error) {
	request := GetEventIndexRequest{
		MsgType:      0xb4,
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetEventIndexResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb4 {
		return nil, errors.New(fmt.Sprintf("GetEventIndex returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.EventIndex{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
	}, nil
}
