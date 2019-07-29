package simulator

import (
	"uhppote/messages"
)

func (s *Simulator) GetDoorDelay(request *messages.GetDoorDelayRequest) (*messages.GetDoorDelayResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	if request.Door < 1 || request.Door > 4 {
		return nil, nil
	}

	response := messages.GetDoorDelayResponse{
		SerialNumber: s.SerialNumber,
		Door:         request.Door,
		Unit:         0x03,
		Delay:        s.Doors[request.Door].Delay.Seconds(),
	}

	return &response, nil
}
