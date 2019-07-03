package simulator

import (
	"time"
	"uhppote"
	"uhppote/types"
)

func (s *Simulator) Find(request *uhppote.FindDevicesRequest) (*uhppote.FindDevicesResponse, error) {
	utc := time.Now().UTC()
	datetime := utc.Add(time.Duration(s.TimeOffset))

	response := uhppote.FindDevicesResponse{
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   s.MacAddress,
		Version:      types.Version(s.Version),
		Date:         types.Date(datetime),
	}

	return &response, nil
}
