package simulator

import (
	"uhppote"
	"uhppote/types"
)

func (s *Simulator) Find(bytes []byte) (*uhppote.FindDevicesResponse, error) {
	response := uhppote.FindDevicesResponse{
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   s.MacAddress,
		Version:      types.Version(s.Version),
		Date:         s.Date,
	}

	return &response, nil
}
