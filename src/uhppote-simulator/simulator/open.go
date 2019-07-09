package simulator

import (
	"uhppote"
)

func (s *Simulator) OpenDoor(request *uhppote.OpenDoorRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	succeeded := request.Door > 0 && request.Door <= 4

	response := uhppote.OpenDoorResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    succeeded,
	}

	return &response, nil
}
