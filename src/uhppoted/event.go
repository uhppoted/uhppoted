package uhppoted

import (
	"context"
	"fmt"
	"os"
	"uhppote"
	"uhppote/types"
)

type Event struct {
	LastEventIndex uint32         `json:"last-event-index"`
	EventType      byte           `json:"event-type"`
	Granted        bool           `json:"access-granted"`
	Door           byte           `json:"door"`
	DoorOpened     bool           `json:"door-opened"`
	UserId         uint32         `json:"user-id"`
	EventTimestamp types.DateTime `json:"event-timestamp"`
	EventResult    byte           `json:"event-result"`
	DoorState      []bool         `json:"door-states"`
	DoorButton     []bool         `json:"door-buttons"`
	SystemState    byte           `json:"system-state"`
	SystemDateTime types.DateTime `json:"system-datetime"`
	PacketNumber   uint32         `json:"packet-number"`
	Backup         uint32         `json:"backup-state"`
	SpecialMessage byte           `json:"special-message"`
	Battery        byte           `json:"battery-status"`
	FireAlarm      byte           `json:"fire-alarm-status"`
}

type EventMessage struct {
	Device struct {
		ID    uint32 `json:"id"`
		Event Event  `json:"event"`
	} `json:"device"`
}

func (u *UHPPOTED) Listen(ctx context.Context, q chan os.Signal) {
	p := make(chan *types.Status)

	go func() {
		if err := ctx.Value("uhppote").(*uhppote.UHPPOTE).Listen(p, q); err != nil {
			u.warn(ctx, 0, "listen", err)
		}
	}()

	for {
		event := <-p
		if event == nil {
			break
		}

		u.log(ctx, "EVENT", uint32(event.SerialNumber), fmt.Sprintf("%v", event))

		message := EventMessage{
			struct {
				ID    uint32 `json:"id"`
				Event Event  `json:"event"`
			}{
				ID: uint32(event.SerialNumber),
				Event: Event{
					LastEventIndex: event.LastIndex,
					EventType:      event.EventType,
					Granted:        event.Granted,
					Door:           event.Door,
					DoorOpened:     event.DoorOpened,
					UserId:         event.UserId,
					EventTimestamp: event.EventTimestamp,
					EventResult:    event.EventResult,
					DoorState:      event.DoorState,
					DoorButton:     event.DoorButton,
					SystemState:    event.SystemState,
					SystemDateTime: event.SystemDateTime,
					PacketNumber:   event.PacketNumber,
					Backup:         event.Backup,
					SpecialMessage: event.SpecialMessage,
					Battery:        event.Battery,
					FireAlarm:      event.FireAlarm,
				},
			},
		}

		u.send(ctx, message)
	}
}
