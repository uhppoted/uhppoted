package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type GetCardsRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x58"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetCardsResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x58"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Records      uint32             `uhppote:"offset:8"`
}
