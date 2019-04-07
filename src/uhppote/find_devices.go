package uhppote

import (
	"errors"
	"fmt"
	"net"
	"uhppote/types"
)

type FindDevicesRequest struct {
	MsgType byte `uhppote:"offset:1"`
}

type FindDevicesResponse struct {
	MsgType      byte             `uhppote:"offset:1"`
	SerialNumber uint32           `uhppote:"offset:4"`
	IpAddress    net.IP           `uhppote:"offset:8"`
	SubnetMask   net.IP           `uhppote:"offset:12"`
	Gateway      net.IP           `uhppote:"offset:16"`
	MacAddress   net.HardwareAddr `uhppote:"offset:20"`
	Version      types.Version    `uhppote:"offset:26"`
	Date         types.Date       `uhppote:"offset:28"`
}

func (u *UHPPOTE) FindDevices() ([]types.Device, error) {
	request := FindDevicesRequest{
		MsgType: 0x94,
	}

	reply := FindDevicesResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x94 {
		return nil, errors.New(fmt.Sprintf("FindDevices returned incorrect message type: %02x\n", reply.MsgType))
	}

	devices := []types.Device{}

	devices = append(devices, types.Device{
		SerialNumber: reply.SerialNumber,
		IpAddress:    reply.IpAddress,
		SubnetMask:   reply.SubnetMask,
		Gateway:      reply.Gateway,
		MacAddress:   reply.MacAddress,
		Version:      reply.Version,
		Date:         reply.Date,
	})

	return devices, nil
}
