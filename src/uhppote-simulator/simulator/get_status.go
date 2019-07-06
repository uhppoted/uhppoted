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
	}

	if event != nil {
		response.SwipeRecord = event.RecordNumber
		response.Granted = event.Granted
		response.Door = event.Door
		response.DoorOpen = event.Opened
		response.CardNumber = event.CardNumber
		response.SwipeDateTime = event.Timestamp
		response.SwipeReason = event.Reason
		response.Door1State = event.Door1State
		response.Door2State = event.Door2State
		response.Door3State = event.Door3State
		response.Door4State = event.Door4State
		response.Door1Button = event.Door1Button
		response.Door2Button = event.Door2Button
		response.Door3Button = event.Door3Button
		response.Door4Button = event.Door4Button
	}

	return &response, nil
}
