package uhppote

import (
	"errors"
	"fmt"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetDoorDelay(serialNumber uint32, door uint8, delay uint8) (*types.DoorDelay, error) {
	request := messages.SetDoorDelayRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
		Unit:         0x03,
		Delay:        delay,
	}

	reply := messages.SetDoorDelayResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.Unit != 0x03 {
		return nil, errors.New(fmt.Sprintf("SetDoorDelay returned incorrect time unit: %02X\n", reply.Unit))
	}

	return &types.DoorDelay{
		SerialNumber: reply.SerialNumber,
		Door:         reply.Door,
		Delay:        reply.Delay,
	}, nil
}
