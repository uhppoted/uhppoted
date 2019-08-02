package simulator

import (
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) SetTime(request *messages.SetTimeRequest) *messages.SetTimeResponse {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	dt := time.Time(request.DateTime).Format("2006-01-02 15:04:05")
	utc, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.UTC)
	if err != nil {
		s.onError(err)
		return nil
	}

	now := time.Now().UTC()
	delta := utc.Sub(now)
	datetime := now.Add(delta)

	s.TimeOffset = Offset(delta)
	err = s.Save()
	if err != nil {
		s.onError(err)
		return nil
	}

	response := messages.SetTimeResponse{
		SerialNumber: s.SerialNumber,
		DateTime:     types.DateTime(datetime),
	}

	return &response
}
