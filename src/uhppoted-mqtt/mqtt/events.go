package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"uhppote/types"
	"uhppoted"
)

type startdate time.Time
type enddate time.Time

func (m *MQTTD) getEvents(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID *uint32    `json:"device-id"`
		Start    *startdate `json:"start"`
		End      *enddate   `json:"end"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "get-events", "Missing/invalid device ID", uhppoted.StatusBadRequest, err)
	} else if body.DeviceID == nil {
		m.OnError(ctx, "get-events", "Missing/invalid device ID", uhppoted.StatusBadRequest, errors.New("Missing device/invalid ID"))
	} else if *body.DeviceID == 0 {
		m.OnError(ctx, "get-events", "Missing/invalid device ID", uhppoted.StatusBadRequest, errors.New("Missing device/invalid ID"))
	} else if body.Start != nil && body.End != nil && time.Time(*body.End).Before(time.Time(*body.Start)) {
		m.OnError(ctx, "get-events", "Invalid date range", uhppoted.StatusBadRequest, fmt.Errorf("Invalid date range (%s to %s)", (*types.DateTime)(body.Start), (*types.DateTime)(body.End)))
	} else {
		rq := uhppoted.GetEventsRequest{
			DeviceID: *body.DeviceID,
			Start:    (*types.DateTime)(body.Start),
			End:      (*types.DateTime)(body.End),
		}

		if response, err := impl.GetEvents(ctx, rq); err != nil {
			m.OnError(ctx, "get-events", "Error retrieving events", uhppoted.StatusInternalServerError, err)
		} else {
			m.Reply(ctx, response)
		}
	}
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
