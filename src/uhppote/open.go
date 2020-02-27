package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

func (u *UHPPOTE) OpenDoor(serialNumber uint32, door uint8) (*types.Result, error) {
	request := messages.OpenDoorRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
	}

	reply := messages.OpenDoorResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
