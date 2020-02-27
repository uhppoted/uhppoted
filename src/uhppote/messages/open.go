package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

type OpenDoorRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x40"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Door         uint8              `uhppote:"offset:8"`
}

type OpenDoorResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x40"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}
