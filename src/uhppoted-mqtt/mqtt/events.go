package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted/src/uhppoted"
	"time"
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Start != nil && body.End != nil && time.Time(*body.End).Before(time.Time(*body.Start)) {
		return nil, ferror(fmt.Errorf("Invalid event date range: %v to %v", body.Start, body.End), "Missing event date range")
	}

	rq := uhppoted.GetEventRangeRequest{
		DeviceID: *body.DeviceID,
		Start:    (*types.DateTime)(body.Start),
		End:      (*types.DateTime)(body.End),
	}

	response, err := impl.GetEventRange(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving events from %v", *body.DeviceID))
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetEventRangeResponse
	}{
		metainfo:              meta,
		GetEventRangeResponse: *response,
	}, nil
}

func (m *MQTTD) getEvent(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		EventID  *uint32            `json:"event-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.EventID == nil {
		return nil, InvalidEventID
	}

	rq := uhppoted.GetEventRequest{
		DeviceID: *body.DeviceID,
		EventID:  *body.EventID,
	}

	response, err := impl.GetEvent(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving events from %v", *body.DeviceID))
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
