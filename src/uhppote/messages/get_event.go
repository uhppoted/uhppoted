package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type GetEventRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb0"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

type GetEventResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb0"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	Type         uint8              `uhppote:"offset:12"`
	Granted      bool               `uhppote:"offset:13"`
	Door         uint8              `uhppote:"offset:14"`
	DoorOpened   bool               `uhppote:"offset:15"`
	UserID       uint32             `uhppote:"offset:16"`
	Timestamp    types.DateTime     `uhppote:"offset:20"`
	Result       uint8              `uhppote:"offset:27"`
}
