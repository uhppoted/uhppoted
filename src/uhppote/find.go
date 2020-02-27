package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

func (u *UHPPOTE) FindDevices() ([]types.Device, error) {
	request := messages.FindDevicesRequest{}
	replies := []messages.FindDevicesResponse{}

	if err := u.Broadcast(request, &replies); err != nil {
		return nil, err
	}

	devices := []types.Device{}
	for _, reply := range replies {
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

	return devices, nil
}

func (u *UHPPOTE) FindDevice(serialNumber uint32) (*types.Device, error) {
	request := messages.FindDevicesRequest{}
	replies := []messages.FindDevicesResponse{}

	if err := u.DirectedBroadcast(serialNumber, request, &replies); err != nil {
		return nil, err
	}

	for _, reply := range replies {
		if uint32(reply.SerialNumber) == serialNumber {
			return &types.Device{
				SerialNumber: reply.SerialNumber,
				IpAddress:    reply.IpAddress,
				SubnetMask:   reply.SubnetMask,
				Gateway:      reply.Gateway,
				MacAddress:   reply.MacAddress,
				Version:      reply.Version,
				Date:         reply.Date,
			}, nil
		}
	}

	return nil, nil
}
