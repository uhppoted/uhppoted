package uhppoted

import (
	"context"
	"fmt"
	"time"
	"uhppote"
	"uhppote/types"
)

type GetEventsRequest struct {
	DeviceID uint32
	Start    *types.DateTime
	End      *types.DateTime
}

type GetEventsResponse struct {
	Device struct {
		ID     uint32      `json:"id"`
		Dates  *DateRange  `json:"dates,omitempty"`
		Events *EventRange `json:"events,omitempty"`
	} `json:"device"`
}

type GetEventRequest struct {
	DeviceID uint32
	EventID  uint32
}

type GetEventResponse struct {
	Device struct {
		ID    uint32 `json:"id"`
		Event event  `json:"event"`
	} `json:"device"`
}

type DateRange struct {
	Start *types.DateTime `json:"start,omitempty"`
	End   *types.DateTime `json:"end,omitempty"`
}

type EventRange struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
}

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

func (u *UHPPOTED) GetEvents(ctx context.Context, request GetEventsRequest) (*GetEventsResponse, int, error) {
	u.debug(ctx, "get-events", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	start := request.Start
	end := request.End

	event, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, 0xffffffff)
	if err != nil {
		return nil, StatusInternalServerError, err
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
				record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, index)
				if err != nil {
					return nil, StatusInternalServerError, err
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

	dates := (*DateRange)(nil)
	if start != nil || end != nil {
		dates = &DateRange{
			Start: start,
			End:   end,
		}
	}

	events := (*EventRange)(nil)
	if first != 0 || last != 0 {
		events = &EventRange{
			First: first,
			Last:  last,
		}
	}

	response := GetEventsResponse{
		Device: struct {
			ID     uint32      `json:"id"`
			Dates  *DateRange  `json:"dates,omitempty"`
			Events *EventRange `json:"events,omitempty"`
		}{
			ID:     device,
			Dates:  dates,
			Events: events,
		},
	}

	u.debug(ctx, "get-events", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

func (u *UHPPOTED) GetEvent(ctx context.Context, request GetEventRequest) (*GetEventResponse, int, error) {
	u.debug(ctx, "get-events", fmt.Sprintf("request  %v", request))

	device := request.DeviceID
	eventID := request.EventID

	record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, eventID)
	if err != nil {
		return nil, StatusInternalServerError, fmt.Errorf("Failed to retrieve event ID %d", eventID)
	}

	if record == nil {
		return nil, StatusNotFound, fmt.Errorf("No event record for ID %d", eventID)
	}

	if record.Index != eventID {
		return nil, StatusNotFound, fmt.Errorf("No event record for ID %d", eventID)
	}

	response := GetEventResponse{
		Device: struct {
			ID    uint32 `json:"id"`
			Event event  `json:"event"`
		}{
			ID: device,
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
		},
	}

	u.debug(ctx, "get-event", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
