package rest

import (
	"context"
	"errors"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"net/http"
)

type event struct {
	Index      uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserID     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

func getEvents(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	first, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, 0)
	if err != nil {
		warn(ctx, deviceID, "get-events", err)
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	if first == nil {
		warn(ctx, deviceID, "get-events", errors.New("No record returned for 'first' event"))
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	last, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, 0xffffffff)
	if err != nil {
		warn(ctx, deviceID, "get-events", err)
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	if last == nil {
		warn(ctx, deviceID, "get-events", errors.New("No record returned for 'last' event"))
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	response := struct {
		Events struct {
			First uint32 `json:"first"`
			Last  uint32 `json:"last"`
		} `json:"events"`
	}{
		Events: struct {
			First uint32 `json:"first"`
			Last  uint32 `json:"last"`
		}{
			First: first.Index,
			Last:  last.Index,
		},
	}

	reply(ctx, w, response)
}

func getEvent(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	eventID := ctx.Value("event-id").(uint32)

	record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, eventID)
	if err != nil {
		warn(ctx, deviceID, "get-event", err)
		http.Error(w, "Error retrieving event", http.StatusInternalServerError)
		return
	}

	if record == nil {
		http.Error(w, "Event record does not exist", http.StatusNotFound)
		return
	}

	if record.Index != eventID {
		http.Error(w, "Event record does not exist", http.StatusNotFound)
		return
	}

	response := struct {
		Event event `json:"event"`
	}{
		Event: event{
			Index:      record.Index,
			Type:       record.Type,
			Granted:    record.Granted,
			Door:       record.Door,
			DoorOpened: record.DoorOpened,
			UserID:     record.UserID,
			Timestamp:  record.Timestamp,
			Result:     record.Result,
		},
	}

	reply(ctx, w, response)
}
