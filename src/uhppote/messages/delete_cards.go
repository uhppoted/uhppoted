package messages

import (
	"uhppote/types"
)

type DeleteCardsRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	MagicNumber  uint32             `uhppote:"offset:8"`
}

type DeleteCardsResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x54"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}