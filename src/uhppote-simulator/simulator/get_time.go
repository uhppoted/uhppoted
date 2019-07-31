package simulator

import (
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) GetTime(request *messages.GetTimeRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	utc := time.Now().UTC()
	datetime := utc.Add(time.Duration(s.TimeOffset))

	response := messages.GetTimeResponse{
		SerialNumber: s.SerialNumber,
		DateTime:     types.DateTime(datetime),
	}

	return &response, nil
}
