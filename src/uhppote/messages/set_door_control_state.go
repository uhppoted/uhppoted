package messages

import (
	"uhppote/types"
)

type SetDoorDelayRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	Unit         uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}

type SetDoorDelayResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	Unit         uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}
