package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type SetDoorControlStateRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	ControlState uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}

type SetDoorControlStateResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x80"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	ControlState uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}
