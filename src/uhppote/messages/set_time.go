package messages

import (
	"uhppote/types"
)

type SetTimeRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x30"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}

type SetTimeResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x30"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	DateTime     types.DateTime     `uhppote:"offset:8"`
}
