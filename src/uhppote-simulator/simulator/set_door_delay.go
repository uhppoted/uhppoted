package simulator

import (
	"uhppote"
)

func (s *Simulator) SetDoorDelay(request *uhppote.SetDoorDelayRequest) (*uhppote.SetDoorDelayResponse, error) {
	if request.SerialNumber != s.SerialNumber {
		return nil, nil
	}

	if request.Unit != 0x03 {
		return nil, nil
	}

	door := request.Door
	if door < 1 || door > 4 {
		return nil, nil
	}

	s.Doors[door].Delay = request.Delay

	err := s.Save()
	if err != nil {
		return nil, err
	}

	response := uhppote.SetDoorDelayResponse{
		SerialNumber: s.SerialNumber,
		Door:         door,
		Unit:         0x03,
		Delay:        s.Doors[door].Delay,
	}

	return &response, nil
}
