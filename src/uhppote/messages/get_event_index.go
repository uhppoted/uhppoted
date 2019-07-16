package messages

import (
	"uhppote/types"
)

type GetEventIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb4"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetEventIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb4"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}
