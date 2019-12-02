package uhppoted

import (
	"context"
	"uhppote"
	"uhppote/types"
)

type GetEventsResponse struct {
	Device struct {
		ID     uint32  `json:"id"`
		Events []event `json:"events"`
	} `json:"device"`
}

type GetEventResponse struct {
	Device struct {
		ID    uint32 `json:"id"`
		Event event  `json:"event"`
	} `json:"device"`
}

type event struct {
	Index      uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserId     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

func (u *UHPPOTED) GetEvents(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-events", rq)

	id, err := rq.DeviceID()
	if err != nil {
		u.warn(ctx, 0, "get-events", err)
		u.oops(ctx, "get-events", "Missing/invalid device ID)", StatusBadRequest)
		return
	}

	last, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(*id, 0xffffffff)
	if err != nil {
		u.warn(ctx, *id, "get-events", err)
		u.oops(ctx, "get-events", "Error retrieving last events", StatusInternalServerError)
		return
	}

	events := make([]event, 0)

	if last != nil {
		for index := last.Index; index > 0; index-- {
			record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(*id, index)
			if err != nil {
				u.warn(ctx, *id, "get-events", err)
				u.oops(ctx, "get-events", "Error retrieving events", StatusInternalServerError)
				return
			}

			events = append(events, event{
				Index:      record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				DoorOpened: record.DoorOpened,
				UserId:     record.UserId,
				Timestamp:  record.Timestamp,
				Result:     record.Result,
			})
		}
	}

	response := GetEventsResponse{
		struct {
			ID     uint32  `json:"id"`
			Events []event `json:"events"`
		}{
			ID:     *id,
			Events: events,
		},
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetEvent(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-event", rq)

	id, eventID, err := rq.DeviceEventID()
	if err != nil {
		u.warn(ctx, 0, "get-event", err)
		u.oops(ctx, "get-event", "Missing/invalid device ID or event ID", StatusBadRequest)
		return
	}

	record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(*id, *eventID)
	if err != nil {
		u.warn(ctx, *id, "get-event", err)
		u.oops(ctx, "get-event", "Failed to retrieve event", StatusInternalServerError)
		return
	}

	if record == nil {
		u.oops(ctx, "get-event", "Event record does not exist", StatusNotFound)
		return
	}

	if record.Index != *eventID {
		u.oops(ctx, "get-event", "Event record does not exist", StatusNotFound)
		return
	}

	response := GetEventResponse{
		struct {
			ID    uint32 `json:"id"`
			Event event  `json:"event"`
		}{
			ID: *id,
			Event: event{
				Index:      record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				DoorOpened: record.DoorOpened,
				UserId:     record.UserId,
				Timestamp:  record.Timestamp,
				Result:     record.Result,
			},
		},
	}

	u.reply(ctx, response)
}
