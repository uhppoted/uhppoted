package uhppote

import (
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetStatus(serialNumber uint32) (*types.Status, error) {
	request := messages.GetStatusRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := messages.GetStatusResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	d := time.Time(reply.SystemDate).Format("2006-01-02")
	t := time.Time(reply.SystemTime).Format("15:04:05")
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", d+" "+t, time.Local)

	return &types.Status{
		SerialNumber:   reply.SerialNumber,
		LastIndex:      reply.LastIndex,
		EventType:      reply.EventType,
		Granted:        reply.Granted,
		Door:           reply.Door,
		DoorOpened:     reply.DoorOpened,
		UserId:         reply.UserId,
		EventTimestamp: reply.EventTimestamp,
		EventResult:    reply.EventResult,
		DoorState:      []bool{reply.Door1State, reply.Door2State, reply.Door3State, reply.Door4State},
		DoorButton:     []bool{reply.Door1Button, reply.Door2Button, reply.Door3Button, reply.Door4Button},
		SystemState:    reply.SystemState,
		SystemDateTime: types.DateTime(datetime),
		PacketNumber:   reply.PacketNumber,
		Backup:         reply.Backup,
		SpecialMessage: reply.SpecialMessage,
		Battery:        reply.Battery,
		FireAlarm:      reply.FireAlarm,
	}, nil
}
