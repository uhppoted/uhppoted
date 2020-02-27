package UT0311L04

import (
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote-simulator/entities"
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
	"time"
)

func (s *UT0311L04) setTime(addr *net.UDPAddr, request *messages.SetTimeRequest) {
	if s.SerialNumber == request.SerialNumber {
		dt := time.Time(request.DateTime).Format("2006-01-02 15:04:05")
		utc, err := time.ParseInLocation("2006-01-02 15:04:05", dt, time.UTC)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		now := time.Now().UTC()
		delta := utc.Sub(now)
		datetime := now.Add(delta)

		s.TimeOffset = entities.Offset(delta)
		response := messages.SetTimeResponse{
			SerialNumber: s.SerialNumber,
			DateTime:     types.DateTime(datetime),
		}

		s.send(addr, &response)

		if err = s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
