package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetDoorDelay(request *messages.GetDoorDelayRequest) *messages.GetDoorDelayResponse {
	if request.SerialNumber != s.SerialNumber {
		return nil
	}

	if request.Door < 1 || request.Door > 4 {
		return nil
	}

	response := messages.GetDoorDelayResponse{
		SerialNumber: s.SerialNumber,
		Door:         request.Door,
		Unit:         0x03,
		Delay:        s.Doors[request.Door].Delay.Seconds(),
	}

	return &response
}
