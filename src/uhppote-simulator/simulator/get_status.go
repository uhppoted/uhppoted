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

	response := uhppote.GetStatusResponse{
		SerialNumber: s.SerialNumber,
		//	LastIndex      uint32             `uhppote:"offset:8"`
		//	SwipeRecord    byte               `uhppote:"offset:12"`
		//	Granted        bool               `uhppote:"offset:13"`
		//	Door           byte               `uhppote:"offset:14"`
		//	DoorOpen       bool               `uhppote:"offset:15"`
		//	CardNumber     uint32             `uhppote:"offset:16"`
		//	SwipeDateTime  types.DateTime     `uhppote:"offset:20"`
		//	SwipeReason    byte               `uhppote:"offset:27"`
		//	Door1State     bool               `uhppote:"offset:28"`
		//	Door2State     bool               `uhppote:"offset:29"`
		//	Door3State     bool               `uhppote:"offset:30"`
		//	Door4State     bool               `uhppote:"offset:31"`
		//	Door1Button    bool               `uhppote:"offset:32"`
		//	Door2Button    bool               `uhppote:"offset:33"`
		//	Door3Button    bool               `uhppote:"offset:34"`
		//	Door4Button    bool               `uhppote:"offset:35"`
		//	SystemState    byte               `uhppote:"offset:36"`
		SystemDate:     types.SystemDate(datetime),
		SystemTime:     types.SystemTime(datetime),
		PacketNumber:   s.PacketNumber,
		Backup:         s.Backup,
		SpecialMessage: s.SpecialMessage,
		Battery:        s.Battery,
		FireAlarm:      s.FireAlarm,
	}

	return &response, nil
}
