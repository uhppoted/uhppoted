package rest

import (
	"context"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net/http"
)

func getStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	status, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetStatus(deviceID)
	if err != nil {
		warn(ctx, deviceID, "get-status", err)
		http.Error(w, "Error retrieving device status", http.StatusInternalServerError)
		return
	}

	response := struct {
		LastEventIndex uint32         `json:"last-event-index"`
		EventType      byte           `json:"event-type"`
		Granted        bool           `json:"access-granted"`
		Door           byte           `json:"door"`
		DoorOpened     bool           `json:"door-opened"`
		UserID         uint32         `json:"user-id"`
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
	}{
		LastEventIndex: status.LastIndex,
		EventType:      status.EventType,
		Granted:        status.Granted,
		Door:           status.Door,
		DoorOpened:     status.DoorOpened,
		UserID:         status.UserID,
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
	}

	reply(ctx, w, response)
}
