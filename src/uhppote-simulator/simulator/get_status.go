package simulator

import (
	"time"
	"uhppote"
	"uhppote/types"
)

func (s *Simulator) GetStatus(request *uhppote.GetStatusRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	utc := time.Now().UTC()
	datetime := utc.Add(time.Duration(s.TimeOffset))
	event := s.Events.Get(s.Events.LastIndex)

	response := uhppote.GetStatusResponse{
		SerialNumber:   s.SerialNumber,
		LastIndex:      s.Events.LastIndex,
		SystemState:    s.SystemState,
		SystemDate:     types.SystemDate(datetime),
		SystemTime:     types.SystemTime(datetime),
		PacketNumber:   s.PacketNumber,
		Backup:         s.Backup,
		SpecialMessage: s.SpecialMessage,
		Battery:        s.Battery,
		FireAlarm:      s.FireAlarm,

		Door1State:  s.Doors[1].Open,
		Door2State:  s.Doors[2].Open,
		Door3State:  s.Doors[3].Open,
		Door4State:  s.Doors[4].Open,
		Door1Button: s.Doors[1].Button,
		Door2Button: s.Doors[2].Button,
		Door3Button: s.Doors[3].Button,
		Door4Button: s.Doors[4].Button,
	}

	if event != nil {
		// response.SwipeRecord = s.Events.LastIndex
		response.Granted = event.Granted
		response.Door = event.Door
		response.DoorOpened = event.DoorOpened
		response.UserId = event.UserId
		response.SwipeDateTime = event.Timestamp
		response.SwipeReason = event.Type
	}

	return &response, nil
}
