package simulator

import (
	"uhppote"
)

func (s *Simulator) GetDoorDelay(request *uhppote.GetDoorDelayRequest) (*uhppote.GetDoorDelayResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	if request.Door < 1 || request.Door > 4 {
		return nil, nil
	}

	response := uhppote.GetDoorDelayResponse{
		SerialNumber: s.SerialNumber,
		Door:         request.Door,
		Unit:         0x03,
		Delay:        s.Doors[request.Door].Delay.Seconds(),
	}

	return &response, nil
}
