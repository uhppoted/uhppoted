package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"uhppote/types"
	"uhppoted"
)

type startdate time.Time
type enddate time.Time

func (m *MQTTD) getEvents(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Start    *startdate         `json:"start"`
		End      *enddate           `json:"end"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	if body.Start != nil && body.End != nil && time.Time(*body.End).Before(time.Time(*body.Start)) {
		return nil, &errorx{
			Err:     fmt.Errorf("Invalid event date range: %v to %v", body.Start, body.End),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing event date range",
		}
	}

	rq := uhppoted.GetEventsRequest{
		DeviceID: *body.DeviceID,
		Start:    (*types.DateTime)(body.Start),
		End:      (*types.DateTime)(body.End),
	}

	response, status, err := impl.GetEvents(rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error retrieving events from %v", *body.DeviceID),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetEventsResponse
	}{
		metainfo:          meta,
		GetEventsResponse: *response,
	}, nil
}

func (m *MQTTD) getEvent(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		EventID  *uint32            `json:"event-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	if body.EventID == nil || *body.EventID == 0 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid event ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid event ID",
		}
	}

	rq := uhppoted.GetEventRequest{
		DeviceID: *body.DeviceID,
		EventID:  *body.EventID,
	}

	response, status, err := impl.GetEvent(rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error retrieving events from %v", *body.DeviceID),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetEventResponse
	}{
		metainfo:         meta,
		GetEventResponse: *response,
	}, nil
}

func (d *startdate) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		*d = startdate(datetime)
		return nil
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		*d = startdate(datetime)
		return nil
	}

	if date, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		*d = startdate(date)
		return nil
	}

	return fmt.Errorf("Cannot parse date/time %s", string(bytes))
}

func (d *enddate) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		*d = enddate(datetime)
		return nil
	}

	if datetime, err := time.ParseInLocation("2006-01-02 15:04", s, time.Local); err == nil {
		*d = enddate(datetime)
		return nil
	}

	if date, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		*d = enddate(time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.Local))
		return nil
	}

	return fmt.Errorf("Cannot parse date/time %s", string(bytes))
}
