package uhppote

import (
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) FindDevices() ([]types.Device, error) {
	request := messages.FindDevicesRequest{}

	replies, err := u.Broadcast(request)
	if err != nil {
		return nil, err
	}

	devices := []types.Device{}
	for _, r := range replies {
		reply := messages.FindDevicesResponse{}
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
