package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
)

type GetListenerResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x92"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Port         uint16             `uhppote:"offset:12"`
}

func (u *UHPPOTE) GetListener(serialNumber uint32) (*types.Listener, error) {
	request := messages.GetListenerRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := messages.GetListenerResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Listener{
		SerialNumber: reply.SerialNumber,
		Address:      net.UDPAddr{IP: reply.Address, Port: int(reply.Port)},
	}, nil
}
