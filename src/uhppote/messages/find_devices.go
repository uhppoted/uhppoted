package messages

import (
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
)

type FindDevicesRequest struct {
	MsgType types.MsgType `uhppote:"value:0x94"`
}

type FindDevicesResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x94"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	IpAddress    net.IP             `uhppote:"offset:8"`
	SubnetMask   net.IP             `uhppote:"offset:12"`
	Gateway      net.IP             `uhppote:"offset:16"`
	MacAddress   types.MacAddress   `uhppote:"offset:20"`
	Version      types.Version      `uhppote:"offset:26"`
	Date         types.Date         `uhppote:"offset:28"`
}
