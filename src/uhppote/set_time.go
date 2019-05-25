package uhppote

import (
	"time"
	"uhppote/types"
)

type SetTimeRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x30"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}

type SetTimeResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x30"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error) {
	request := SetTimeRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		DateTime:     types.DateTime(datetime),
	}

	reply := SetTimeResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Time{
		SerialNumber: reply.SerialNumber,
		DateTime:     reply.DateTime,
	}, nil
}
