package UT0311L04

import (
	"github.com/uhppoted/uhppote-core/messages"
	"github.com/uhppoted/uhppote-core/types"
	"net"
	"time"
)

func (s *UT0311L04) getTime(addr *net.UDPAddr, request *messages.GetTimeRequest) {
	if s.SerialNumber == request.SerialNumber {

		utc := time.Now().UTC()
		datetime := utc.Add(time.Duration(s.TimeOffset))

		response := messages.GetTimeResponse{
			SerialNumber: s.SerialNumber,
			DateTime:     types.DateTime(datetime),
		}

		s.send(addr, &response)
	}
}
