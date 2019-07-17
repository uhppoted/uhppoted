package messages

import (
	"uhppote/types"
)

type SetEventIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb2"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	MagicWord    uint32             `uhppote:"offset:12"`
}

type SetEventIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb2"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Changed      bool               `uhppote:"offset:8"`
}
