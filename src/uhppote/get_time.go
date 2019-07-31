package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetTime(serialNumber uint32) (*types.Time, error) {
	request := messages.GetTimeRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := messages.GetTimeResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Time{
		SerialNumber: reply.SerialNumber,
		DateTime:     reply.DateTime,
	}, nil
}
