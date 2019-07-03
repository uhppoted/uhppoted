package simulator

import (
	"time"
	"uhppote"
	"uhppote/types"
)

func (s *Simulator) GetTime(request *uhppote.GetTimeRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	utc := time.Now().UTC()
	datetime := utc.Add(time.Duration(s.TimeOffset))

	response := uhppote.GetTimeResponse{
		SerialNumber: s.SerialNumber,
		DateTime:     types.DateTime(datetime),
	}

	return &response, nil
}
