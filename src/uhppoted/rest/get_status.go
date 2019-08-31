package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"uhppote"
	"uhppote/types"
)

type Status struct {
	SerialNumber   uint32         `json:"serial-number"`
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

func getStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	matches := regexp.MustCompile("^/uhppote/device/([0-9]+)/status$").FindStringSubmatch(url)
	deviceId, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	status, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetStatus(uint32(deviceId))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device status: %v", err), http.StatusInternalServerError)
		return
	}

	response := Status{
		SerialNumber:   uint32(status.SerialNumber),
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
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
