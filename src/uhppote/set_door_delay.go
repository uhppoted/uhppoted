package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type SetDoorDelayRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	Unit         uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}

type SetDoorDelayResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	Unit         uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}

func (u *UHPPOTE) SetDoorDelay(serialNumber uint32, door uint8, delay uint8) (*types.DoorDelay, error) {
	request := SetDoorDelayRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
		Unit:         0x03,
		Delay:        delay,
	}

	reply := SetDoorDelayResponse{}

	err := u.Execute(request, &reply)
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
