package uhppote

import (
	"errors"
	"fmt"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetDoorDelay(serialNumber uint32, door byte) (*types.DoorDelay, error) {
	request := messages.GetDoorDelayRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
	}

	reply := messages.GetDoorDelayResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.Unit != 0x03 {
		return nil, errors.New(fmt.Sprintf("GetDoorDelay returned incorrect time unit: %02X\n", reply.Unit))
	}

	return &types.DoorDelay{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		Delay:        reply.Delay,
	}, nil
}
