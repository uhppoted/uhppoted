package UT0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UT0311L04) getEvent(addr *net.UDPAddr, request *messages.GetEventRequest) {
	if s.SerialNumber == request.SerialNumber {
		index := request.Index

		if event := s.Events.Get(index); event != nil {
			response := messages.GetEventResponse{
				SerialNumber: s.SerialNumber,
				Index:        event.RecordNumber,
				Type:         event.Type,
				Granted:      event.Granted,
				Door:         event.Door,
				DoorOpened:   event.DoorOpened,
				UserID:       event.UserID,
				Timestamp:    event.Timestamp,
				Result:       event.Result,
			}

			s.send(addr, &response)
		}
	}
}
