package simulator

import (
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) find(request *messages.FindDevicesRequest) *messages.FindDevicesResponse {
	utc := time.Now().UTC()
	datetime := utc.Add(time.Duration(s.TimeOffset))

	response := messages.FindDevicesResponse{
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   s.MacAddress,
		Version:      types.Version(s.Version),
		Date:         types.Date(datetime),
	}

	return &response
}
