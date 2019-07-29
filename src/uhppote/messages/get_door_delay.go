package messages

import (
	"uhppote/types"
)

type GetDoorDelayRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x82"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
}

type GetDoorDelayResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x82"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	Unit         uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}
