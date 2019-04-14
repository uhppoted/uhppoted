package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type SetSwipeIndexRequest struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	MagicWord    uint32             `uhppote:"offset:12"`
}

type SetSwipeIndexResponse struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Success      bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetSwipeIndex(serialNumber, index uint32) (*types.SwipeIndexResult, error) {
	request := SetSwipeIndexRequest{
		MsgType:      0xb2,
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
		MagicWord:    0x55aaaa55,
	}

	reply := SetSwipeIndexResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb2 {
		return nil, errors.New(fmt.Sprintf("GetSwipeIndex returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.SwipeIndexResult{
		SerialNumber: reply.SerialNumber,
		Index:        index,
		Succeeded:    reply.Success,
	}, nil
}
