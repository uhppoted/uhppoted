package uhppote

import (
	"uhppote/types"
)

type OpenDoorRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x40"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
}

type OpenDoorResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x40"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Success      bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) OpenDoor(serialNumber uint32, door uint8) (*types.Result, error) {
	request := OpenDoorRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Door:         door,
	}

	reply := OpenDoorResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Success:      reply.Success,
	}, nil
}
