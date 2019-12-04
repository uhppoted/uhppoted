package uhppoted

import (
	"context"
	"time"
	"uhppote"
	"uhppote/types"
)

type GetEventsResponse struct {
	Device struct {
		ID     uint32      `json:"id"`
		Dates  *daterange  `json:"dates,omitempty"`
		Events *eventrange `json:"events,omitempty"`
	} `json:"device"`
}

type GetEventResponse struct {
	Device struct {
		ID    uint32 `json:"id"`
		Event event  `json:"event"`
	} `json:"device"`
}

type daterange struct {
	Start *types.DateTime `json:"start,omitempty"`
	End   *types.DateTime `json:"end,omitempty"`
}

type eventrange struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
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

	start, end, err := rq.DateRange()
	if err != nil {
		u.warn(ctx, 0, "get-events", err)
		u.oops(ctx, "get-events", "Invalid date range)", StatusBadRequest)
		return
	}

	event, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(*id, 0xffffffff)
	if err != nil {
		u.warn(ctx, *id, "get-events", err)
		u.oops(ctx, "get-events", "Error retrieving last events", StatusInternalServerError)
		return
	}

	first := uint32(0)
	last := uint32(0)
	if event != nil {
		first = 1
		last = event.Index

		if start != nil {
			first = last
		}

		if end != nil {
			last = 1
		}

		if start != nil || end != nil {
			for index := event.Index; index > 0; index-- {
				record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(*id, index)
				if err != nil {
					u.warn(ctx, *id, "get-events", err)
					u.oops(ctx, "get-events", "Error retrieving events", StatusInternalServerError)
					return
				}

				if start != nil && !time.Time(record.Timestamp).Before(time.Time(*start)) && record.Index < first {
					first = record.Index
				}

				if end != nil && !time.Time(*end).Before(time.Time(record.Timestamp)) && record.Index > last {
					last = record.Index
				}
			}
		}
	}

	dates := (*daterange)(nil)
	if start != nil || end != nil {
		dates = &daterange{
			Start: start,
			End:   end,
		}
	}

	events := (*eventrange)(nil)
	if first != 0 || last != 0 {
		events = &eventrange{
			First: first,
			Last:  last,
		}
	}

	response := GetEventsResponse{
		struct {
			ID     uint32      `json:"id"`
			Dates  *daterange  `json:"dates,omitempty"`
			Events *eventrange `json:"events,omitempty"`
		}{
			ID:     *id,
			Dates:  dates,
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
