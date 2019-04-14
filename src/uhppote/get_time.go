package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetTimeRequest struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetTimeResponse struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetTime(serialNumber uint32) (*types.Time, error) {
	request := GetTimeRequest{
		MsgType:      0x32,
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetTimeResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x32 {
		return nil, errors.New(fmt.Sprintf("GetTime returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.Time{
		SerialNumber: reply.SerialNumber,
		DateTime:     reply.DateTime,
	}, nil
}
