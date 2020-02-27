package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type DeleteCardsRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	MagicWord    uint32             `uhppote:"offset:8"`
}

type DeleteCardsResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}
