package simulator

import (
	"uhppote"
)

func (s *Simulator) GetDoorDelay(request *uhppote.GetDoorDelayRequest) (*uhppote.GetDoorDelayResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	door := request.Door

	if door > 0 && door <= 4 {
		response := uhppote.GetDoorDelayResponse{
			SerialNumber: s.SerialNumber,
			Door:         door,
			Unit:         0x03,
			Delay:        s.Doors[door].Delay,
		}

		return &response, nil
	}

	return nil, nil
}
