package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
)

type SetListenerRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x90"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Port         uint16             `uhppote:"offset:12"`
}

type SetListenerResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x90"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}
