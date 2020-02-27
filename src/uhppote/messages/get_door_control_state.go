package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type GetDoorControlStateRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x82"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
}

type GetDoorControlStateResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x82"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
	ControlState uint8              `uhppote:"offset:9"`
	Delay        uint8              `uhppote:"offset:10"`
}
