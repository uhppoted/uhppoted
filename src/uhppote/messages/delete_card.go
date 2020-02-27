package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type DeleteCardRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x52"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
}

type DeleteCardResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x52"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}
