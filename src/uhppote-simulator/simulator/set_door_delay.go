package simulator

import (
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
)

func (s *Simulator) setDoorDelay(request *messages.SetDoorDelayRequest) *messages.SetDoorDelayResponse {
	if request.SerialNumber != s.SerialNumber {
		return nil
	}

	if request.Unit != 0x03 {
		return nil
	}

	door := request.Door
	if door < 1 || door > 4 {
		return nil
	}

	s.Doors[door].Delay = entities.Delay(uint64(request.Delay) * 1000000000)

	err := s.Save()
	if err != nil {
		s.onError(err)
	}

	response := messages.SetDoorDelayResponse{
		SerialNumber: s.SerialNumber,
		Door:         door,
		Unit:         0x03,
		Delay:        s.Doors[door].Delay.Seconds(),
	}

	return &response
}
