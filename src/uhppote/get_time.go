package uhppote

import (
	"uhppote/types"
)

type GetTimeRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x32"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetTimeResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x32"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetTime(serialNumber uint32) (*types.Time, error) {
	request := GetTimeRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetTimeResponse{}

	err := u.Execute(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Time{
		SerialNumber: reply.SerialNumber,
		DateTime:     reply.DateTime,
	}, nil
}
