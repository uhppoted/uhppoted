package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"time"
)

func (u *UHPPOTE) SetTime(serialNumber uint32, datetime time.Time) (*types.Time, error) {
	request := messages.SetTimeRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		DateTime:     types.DateTime(datetime),
	}

	reply := messages.SetTimeResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Time{
		SerialNumber: reply.SerialNumber,
		DateTime:     reply.DateTime,
	}, nil
}
