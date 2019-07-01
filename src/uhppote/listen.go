package uhppote

import (
	"os"
	"time"
	"uhppote/types"
)

type Event struct {
	MsgType        types.MsgType      `uhppote:"value:0x20"`
	SerialNumber   types.SerialNumber `uhppote:"offset:4"`
	LastIndex      uint32             `uhppote:"offset:8"`
	SwipeRecord    byte               `uhppote:"offset:12"`
	Granted        bool               `uhppote:"offset:13"`
	Door           byte               `uhppote:"offset:14"`
	DoorOpen       bool               `uhppote:"offset:15"`
	CardNumber     uint32             `uhppote:"offset:16"`
	SwipeDateTime  types.DateTime     `uhppote:"offset:20"`
	SwipeReason    byte               `uhppote:"offset:27"`
	Door1State     bool               `uhppote:"offset:28"`
	Door2State     bool               `uhppote:"offset:29"`
	Door3State     bool               `uhppote:"offset:30"`
	Door4State     bool               `uhppote:"offset:31"`
	Door1Button    bool               `uhppote:"offset:32"`
	Door2Button    bool               `uhppote:"offset:33"`
	Door3Button    bool               `uhppote:"offset:34"`
	Door4Button    bool               `uhppote:"offset:35"`
	SystemState    byte               `uhppote:"offset:36"`
	SystemDate     types.SystemDate   `uhppote:"offset:51"`
	SystemTime     types.SystemTime   `uhppote:"offset:37"`
	PacketNumber   uint32             `uhppote:"offset:40"` // TODO verify
	Backup         uint32             `uhppote:"offset:44"` // TODO verify
	SpecialMessage byte               `uhppote:"offset:48"` // TODO verify
	LowBattery     byte               `uhppote:"offset:49"` // TODO verify
	FireAlarm      byte               `uhppote:"offset:50"` // TODO verify
}

func (u *UHPPOTE) Listen(p chan *types.Status, q chan os.Signal) error {
	pipe := make(chan Event)

	go func() {
		for {
			event := <-pipe
			p <- event.transform()
		}
	}()

	return u.listen(pipe, q)
}

func (event Event) transform() *types.Status {
	d := time.Time(event.SystemDate).Format("2006-01-02")
	t := time.Time(event.SystemTime).Format("15:04:05")
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", d+" "+t, time.Local)

	return &types.Status{
		SerialNumber:   event.SerialNumber,
		LastIndex:      event.LastIndex,
		SwipeRecord:    event.SwipeRecord,
		Granted:        event.Granted,
		Door:           event.Door,
		DoorOpen:       event.DoorOpen,
		CardNumber:     event.CardNumber,
		SwipeDateTime:  event.SwipeDateTime,
		SwipeReason:    event.SwipeReason,
		DoorState:      []bool{event.Door1State, event.Door2State, event.Door3State, event.Door4State},
		DoorButton:     []bool{event.Door1Button, event.Door2Button, event.Door3Button, event.Door4Button},
		SystemState:    event.SystemState,
		SystemDateTime: types.DateTime(datetime),
		PacketNumber:   event.PacketNumber,
		Backup:         event.Backup,
		SpecialMessage: event.SpecialMessage,
		//LowBattery:     event.LowBattery,
		FireAlarm: event.FireAlarm,
	}
}
