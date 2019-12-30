package mqtt

import (
	"context"
	"encoding/json"
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
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
	} else if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if *body.DeviceID == 0 {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if body.Start != nil && body.End != nil && time.Time(*body.End).Before(time.Time(*body.Start)) {
		m.OnError(ctx, "Invalid date range", uhppoted.StatusBadRequest, fmt.Errorf("Invalid date range '%s'", string(msg.Payload())))
	} else {
		rq := uhppoted.GetEventsRequest{
			DeviceID: *body.DeviceID,
			Start:    (*types.DateTime)(body.Start),
			End:      (*types.DateTime)(body.End),
		}

		if response, status, err := impl.GetEvents(ctx, rq); err != nil {
			m.OnError(ctx, "Error retrieving events", status, err)
		} else if response != nil {
			m.Reply(ctx, response)
		}
	}
}

func (m *MQTTD) getEvent(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		EventID  *uint32 `json:"event-id"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
	} else if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if *body.DeviceID == 0 {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
	} else if body.EventID == nil {
		m.OnError(ctx, "Missing/invalid event ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid event ID '%s'", string(msg.Payload())))
	} else if *body.EventID == 0 {
		m.OnError(ctx, "Missing/invalid event ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid event ID '%s'", string(msg.Payload())))
	} else {
		rq := uhppoted.GetEventRequest{
			DeviceID: *body.DeviceID,
			EventID:  *body.EventID,
		}

		if response, status, err := impl.GetEvent(ctx, rq); err != nil {
			m.OnError(ctx, "Error retrieving events", status, err)
		} else if response != nil {
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
