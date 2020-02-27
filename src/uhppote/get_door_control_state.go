package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

func (u *UHPPOTE) GetDoorControlState(serialNumber uint32, door byte) (*types.DoorControlState, error) {
	request := messages.GetDoorControlStateRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
	}

	reply := messages.GetDoorControlStateResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.DoorControlState{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		ControlState: reply.ControlState,
		Delay:        reply.Delay,
	}, nil
}
