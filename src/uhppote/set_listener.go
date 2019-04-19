package uhppote

import (
	"errors"
	"fmt"
	"net"
	"uhppote/types"
)

type SetListenerRequest struct {
	MsgType      byte               `uhppote:"offset:1, value:0x90"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Port         uint16             `uhppote:"offset:12"`
}

type SetListenerResponse struct {
	MsgType      byte               `uhppote:"offset:1, value:0x90"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Success      bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetListener(serialNumber uint32, address net.UDPAddr) (*types.ListenerResult, error) {
	if address.IP.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP address: %v", address))
	}

	request := SetListenerRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Address:      address.IP,
		Port:         uint16(address.Port),
	}

	reply := SetListenerResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.ListenerResult{
		SerialNumber: reply.SerialNumber,
		Address:      address,
		Succeeded:    reply.Success,
	}, nil
}
