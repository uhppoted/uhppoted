package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

/*
Doesn't seem to be supported by the latest UHPPOTE firmware.
*/

type SetSwipeIndexRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Index        uint32 `uhppote:"offset:8"`
}

type SetSwipeIndexResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Index        uint32 `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetSwipeIndex(serialNumber, index uint32) (*types.SwipeIndex, error) {
	request := SetSwipeIndexRequest{
		MsgType:      0xb2,
		SerialNumber: serialNumber,
		Index:        index,
	}

	reply := SetSwipeIndexResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb2 {
		return nil, errors.New(fmt.Sprintf("GStSwipeIndex returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.SwipeIndex{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
	}, nil
}
