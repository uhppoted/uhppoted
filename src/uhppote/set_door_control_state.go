package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetDoorControlState(serialNumber uint32, door uint8, state uint8, delay uint8) (*types.DoorControlState, error) {
	request := messages.SetDoorControlStateRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
		ControlState: state,
		Delay:        delay,
	}

	reply := messages.SetDoorControlStateResponse{}

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
