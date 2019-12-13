package uhppoted

import (
	"context"
	"fmt"
	"os"
	"uhppote"
	"uhppote/types"
)

type ListenEvent struct {
	DeviceID   uint32         `json:"device-id"`
	EventID    uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserId     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

type EventMessage struct {
	Event ListenEvent `json:"event"`
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

		device := uint32(event.SerialNumber)
		eventID := event.LastIndex
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, eventID)
		if err != nil {
			u.warn(ctx, device, "listen", fmt.Errorf("Failed to retrieve event ID %d", eventID))
			continue
		}

		if record == nil {
			u.warn(ctx, device, "listen", fmt.Errorf("No event record for ID %d", eventID))
			continue
		}

		if record.Index != eventID {
			u.warn(ctx, device, "listen", fmt.Errorf("No event record for ID %d", eventID))
			continue
		}

		message := EventMessage{
			Event: ListenEvent{
				DeviceID:   device,
				EventID:    record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				DoorOpened: record.DoorOpened,
				UserId:     record.UserId,
				Timestamp:  record.Timestamp,
				Result:     record.Result,
			},
		}

		u.debug(ctx, "listen", fmt.Sprintf("event %v", message))
		u.send(ctx, message)
	}
}
