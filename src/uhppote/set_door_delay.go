package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetDoorDelay(serialNumber uint32, door uint8, delay uint8) (*types.DoorDelay, error) {
	request := messages.SetDoorControlStateRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
		ControlState: 3,
		Delay:        delay,
	}

	reply := messages.SetDoorControlStateResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.DoorDelay{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		ControlState: reply.ControlState,
		Delay:        reply.Delay,
	}, nil
}
