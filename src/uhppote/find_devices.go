package uhppote

import (
	"net"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
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
	MacAddress   net.HardwareAddr   `uhppote:"offset:20"`
	Version      types.Version      `uhppote:"offset:26"`
	Date         types.Date         `uhppote:"offset:28"`
}

func (u *UHPPOTE) FindDevices() ([]types.Device, error) {
	request := FindDevicesRequest{}

	replies, err := u.Broadcast(request)
	if err != nil {
		return nil, err
	}

	devices := []types.Device{}
	for _, r := range replies {
		reply := FindDevicesResponse{}
		err = codec.Unmarshal(r, &reply)
		if err != nil {
			return devices, err
		} else {
			devices = append(devices, types.Device{
				SerialNumber: reply.SerialNumber,
				IpAddress:    reply.IpAddress,
				SubnetMask:   reply.SubnetMask,
				Gateway:      reply.Gateway,
				MacAddress:   reply.MacAddress,
				Version:      reply.Version,
				Date:         reply.Date,
			})
		}
	}

	return devices, nil
}
