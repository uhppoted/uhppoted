package rest

import (
	"context"
	"net/http"
	"uhppote"
	"uhppote/types"
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

	last, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, 0xffffffff)
	if err != nil {
		warn(ctx, deviceID, "get-events", err)
		http.Error(w, "Error retrieving events", http.StatusInternalServerError)
		return
	}

	events := make([]event, 0)

	if last != nil {
		for index := last.Index; index > 0; index-- {
			record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(deviceID, index)
			if err != nil {
				warn(ctx, deviceID, "get-event", err)
				http.Error(w, "Error retrieving event", http.StatusInternalServerError)
				return
			}

			events = append(events, event{
				Index:      record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				DoorOpened: record.DoorOpened,
				UserID:     record.UserID,
				Timestamp:  record.Timestamp,
				Result:     record.Result,
			})
		}
	}

	response := struct {
		Events []event `json:"events"`
	}{
		Events: events,
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
