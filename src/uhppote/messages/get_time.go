package messages

import (
	"uhppote/types"
)

type GetTimeRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x32"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetTimeResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x32"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}
