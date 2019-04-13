package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetSwipesRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
}

type GetSwipesResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Count        uint32 `uhppote:"offset:16"`
}

func (u *UHPPOTE) GetSwipeCount(serialNumber uint32) (*types.SwipeCount, error) {
	request := GetSwipesRequest{
		MsgType:      0xb4,
		SerialNumber: serialNumber,
	}

	reply := GetSwipesResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb4 {
		return nil, errors.New(fmt.Sprintf("GetSwipeCount returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.SwipeCount{
		SerialNumber: reply.SerialNumber,
		Count:        reply.Count,
	}, nil
}
