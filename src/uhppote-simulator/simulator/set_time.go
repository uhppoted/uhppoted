package simulator

import (
	"net"
	"time"
	"uhppote-simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) setTime(addr *net.UDPAddr, request *messages.SetTimeRequest) {
	if s.SerialNumber == request.SerialNumber {
		dt := time.Time(request.DateTime).Format("2006-01-02 15:04:05")
		utc, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.UTC)
		if err != nil {
			s.onError(err)
			return
		}

		now := time.Now().UTC()
		delta := utc.Sub(now)
		datetime := now.Add(delta)

		s.TimeOffset = entities.Offset(delta)
		err = s.Save()
		if err != nil {
			s.onError(err)
			return
		}

		response := messages.SetTimeResponse{
			SerialNumber: s.SerialNumber,
			DateTime:     types.DateTime(datetime),
		}

		s.send(addr, &response)
	}
}
