package uhppote

import (
	"time"
	"uhppote/messages"
	"uhppote/types"
)

//type GetStatusRequest struct {
//	MsgType      types.MsgType      `uhppote:"value:0x20"`
//	SerialNumber types.SerialNumber `uhppote:"offset:4"`
//}
//
//type GetStatusResponse struct {
//	MsgType        types.MsgType      `uhppote:"value:0x20"`
//	SerialNumber   types.SerialNumber `uhppote:"offset:4"`
//	LastIndex      uint32             `uhppote:"offset:8"`
//	SwipeRecord    byte               `uhppote:"offset:12"`
//	Granted        bool               `uhppote:"offset:13"`
//	Door           byte               `uhppote:"offset:14"`
//	DoorOpened     bool               `uhppote:"offset:15"`
//	UserId         uint32             `uhppote:"offset:16"`
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
//	SystemDate     types.SystemDate   `uhppote:"offset:51"`
//	SystemTime     types.SystemTime   `uhppote:"offset:37"`
//	PacketNumber   uint32             `uhppote:"offset:40"` // TODO unverified - trust at own risk
//	Backup         uint32             `uhppote:"offset:44"` // TODO unverified - trust at own risk
//	SpecialMessage byte               `uhppote:"offset:48"` // TODO unverified - trust at own risk
//	Battery        byte               `uhppote:"offset:49"` // TODO unverified - trust at own risk
//	FireAlarm      byte               `uhppote:"offset:50"` // TODO unverified - trust at own risk
//}

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
		SwipeRecord:    reply.SwipeRecord,
		Granted:        reply.Granted,
		Door:           reply.Door,
		DoorOpened:     reply.DoorOpened,
		UserId:         reply.UserId,
		SwipeDateTime:  reply.SwipeDateTime,
		SwipeReason:    reply.SwipeReason,
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
