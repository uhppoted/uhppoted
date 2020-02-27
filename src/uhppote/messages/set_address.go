package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
)

type SetAddressRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x96"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Mask         net.IP             `uhppote:"offset:12"`
	Gateway      net.IP             `uhppote:"offset:16"`
	MagicWord    uint32             `uhppote:"offset:20"`
}
