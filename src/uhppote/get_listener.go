package uhppote

import (
	"net"
	"uhppote/types"
)

type GetListenerRequest struct {
	MsgType      byte               `uhppote:"offset:1, value:0x92"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetListenerResponse struct {
	MsgType      byte               `uhppote:"offset:1, value:0x92"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Port         uint16             `uhppote:"offset:12"`
}

func (u *UHPPOTE) GetListener(serialNumber uint32) (*types.Listener, error) {
	request := GetListenerRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetListenerResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Listener{
		SerialNumber: reply.SerialNumber,
		Address:      reply.Address,
		Port:         reply.Port,
	}, nil
}
