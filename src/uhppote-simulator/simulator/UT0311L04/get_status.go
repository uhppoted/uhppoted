package UT0311L04

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net"
	"time"
)

func (s *UT0311L04) getStatus(addr *net.UDPAddr, request *messages.GetStatusRequest) {
	if s.SerialNumber == request.SerialNumber {
		utc := time.Now().UTC()
		datetime := utc.Add(time.Duration(s.TimeOffset))
		event := s.Events.Get(s.Events.Last)

		response := messages.GetStatusResponse{
			SerialNumber:   s.SerialNumber,
			LastIndex:      s.Events.Last,
			SystemState:    s.SystemState,
			SystemDate:     types.SystemDate(datetime),
			SystemTime:     types.SystemTime(datetime),
			PacketNumber:   s.PacketNumber,
			Backup:         s.Backup,
			SpecialMessage: s.SpecialMessage,
			Battery:        s.Battery,
			FireAlarm:      s.FireAlarm,

			Door1State: s.Doors[1].IsOpen(),
			Door2State: s.Doors[2].IsOpen(),
			Door3State: s.Doors[3].IsOpen(),
			Door4State: s.Doors[4].IsOpen(),

			Door1Button: s.Doors[1].IsButtonPressed(),
			Door2Button: s.Doors[2].IsButtonPressed(),
			Door3Button: s.Doors[3].IsButtonPressed(),
			Door4Button: s.Doors[4].IsButtonPressed(),
		}

		if event != nil {
			response.EventType = event.Type
			response.Granted = event.Granted
			response.Door = event.Door
			response.DoorOpened = event.DoorOpened
			response.UserID = event.UserID
			response.EventTimestamp = event.Timestamp
			response.EventResult = event.Result
		}

		s.send(addr, &response)
	}
}
