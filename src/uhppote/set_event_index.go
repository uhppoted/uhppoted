package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type SetEventIndexRequest struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	MagicWord    uint32             `uhppote:"offset:12"`
}

type SetEventIndexResponse struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Success      bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetEventIndex(serialNumber, index uint32) (*types.EventIndexResult, error) {
	request := SetEventIndexRequest{
		MsgType:      0xb2,
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
		MagicWord:    0x55aaaa55,
	}

	reply := SetEventIndexResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb2 {
		return nil, errors.New(fmt.Sprintf("GetEventIndex returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.EventIndexResult{
		SerialNumber: reply.SerialNumber,
		Index:        index,
		Succeeded:    reply.Success,
	}, nil
}
