package uhppote

import (
	"errors"
	"fmt"
	"net"
	"uhppote/encoding"
	"uhppote/types"
)

func (u *UHPPOTE) FindDevices() ([]types.Device, error) {
	request := struct {
		MsgType byte `uhppote:"offset:1"`
	}{
		0x94,
	}

	reply := struct {
		MsgType      byte             `uhppote:"offset:1"`
		SerialNumber uint32           `uhppote:"offset:4"`
		IpAddress    net.IP           `uhppote:"offset:8"`
		SubnetMask   net.IP           `uhppote:"offset:12"`
		Gateway      net.IP           `uhppote:"offset:16"`
		MacAddress   net.HardwareAddr `uhppote:"offset:20"`
		Version      types.Version    `uhppote:"offset:26"`
		Date         types.Date       `uhppote:"offset:28"`
	}{}

	p, err := uhppote.Marshal(request)
	if err != nil {
		return nil, err
	}

	q, err := u.Execute(p)
	if err != nil {
		return nil, err
	}

	err = uhppote.Unmarshal(q, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x94 {
		return nil, errors.New(fmt.Sprint("FindDevices returned incorrect message type: %02x\n", reply.MsgType))
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
