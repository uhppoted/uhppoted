package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
)

type GetListenerRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x92"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetListenerResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x92"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Port         uint16             `uhppote:"offset:12"`
}
