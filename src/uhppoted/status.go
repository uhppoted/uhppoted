package uhppoted

import (
	"context"
	"uhppote"
	"uhppote/types"
)

type Status struct {
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

type DeviceStatus struct {
	Device struct {
		ID     uint32 `json:"id"`
		Status Status `json:"status"`
	} `json:"device"`
}

func (u *UHPPOTED) GetStatus(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-status", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, 0, "get-status", err)
		u.oops(ctx, "get-status", "Error retrieving device status (invalid device ID)", StatusBadRequest)
		return
	}

	status, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetStatus(*id)
	if err != nil {
		u.warn(ctx, *id, "get-status", err)
		u.oops(ctx, "get-status", "Error retrieving device status", StatusInternalServerError)
		return
	}

	response := DeviceStatus{
		struct {
			ID     uint32 `json:"id"`
			Status Status `json:"status"`
		}{
			ID: *id,
			Status: Status{
				LastEventIndex: status.LastIndex,
				EventType:      status.EventType,
				Granted:        status.Granted,
				Door:           status.Door,
				DoorOpened:     status.DoorOpened,
				UserId:         status.UserId,
				EventTimestamp: status.EventTimestamp,
				EventResult:    status.EventResult,
				DoorState:      status.DoorState,
				DoorButton:     status.DoorButton,
				SystemState:    status.SystemState,
				SystemDateTime: status.SystemDateTime,
				PacketNumber:   status.PacketNumber,
				Backup:         status.Backup,
				SpecialMessage: status.SpecialMessage,
				Battery:        status.Battery,
				FireAlarm:      status.FireAlarm,
			},
		},
	}

	u.reply(ctx, response)
}
