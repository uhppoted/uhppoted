package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

/*
Doesn't seem to be supported by the latest UHPPOTE firmware - returns 0.

Use either GetSwipe(<serialnumber>, 0xffffffff) to get the latest swipe
or GetStatus(<serialnumber>) to get the controller status which includes
the swipe index.
*/

type GetSwipeIndexRequest struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
}

type GetSwipeIndexResponse struct {
	MsgType      byte   `uhppote:"offset:1"`
	SerialNumber uint32 `uhppote:"offset:4"`
	Index        uint32 `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetSwipeIndex(serialNumber uint32) (*types.SwipeIndex, error) {
	request := GetSwipeIndexRequest{
		MsgType:      0xb4,
		SerialNumber: serialNumber,
	}

	reply := GetSwipeIndexResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb4 {
		return nil, errors.New(fmt.Sprintf("GetSwipeIndex returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.SwipeIndex{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
	}, nil
}
