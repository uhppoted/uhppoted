package messages

import (
	"uhppote/types"
)

type PutCardRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x50"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         types.Date         `uhppote:"offset:12"`
	To           types.Date         `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

type PutCardResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x50"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}
