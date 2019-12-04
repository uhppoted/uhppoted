package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
	"uhppote/types"
)

type Request struct {
	Message MQTT.Message
}

func (rq Request) String() string {
	return rq.Message.Topic() + "  " + string(rq.Message.Payload())
}

func (rq *Request) DeviceID() (*uint32, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, err
	} else if body.DeviceID == nil {
		return nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, errors.New("Missing device ID")
	}

	return body.DeviceID, nil
}

func (rq *Request) DateTime() (*time.Time, error) {
	body := struct {
		DateTime *types.DateTime `json:"datetime"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, err
	} else if body.DateTime == nil {
		return nil, errors.New("Missing date/time")
	}

	return (*time.Time)(body.DateTime), nil
}

func (rq *Request) DeviceDoor() (*uint32, *uint8, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		Door     *uint8  `json:"door"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, errors.New("Missing device ID")
	}

	if body.Door == nil {
		return nil, nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, nil, errors.New("Invalid door")
	}

	return body.DeviceID, body.Door, nil
}

func (rq *Request) DeviceDoorDelay() (*uint32, *uint8, *uint8, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		Door     *uint8  `json:"door"`
		Delay    *uint8  `json:"delay"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, nil, errors.New("Missing device ID")
	}

	if body.Door == nil {
		return nil, nil, nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, nil, nil, errors.New("Invalid door")
	}

	if body.Delay == nil {
		return nil, nil, nil, errors.New("Invalid door delay")
	} else if *body.Delay == 0 || *body.Delay > 60 {
		return nil, nil, nil, errors.New("Invalid door delay")
	}

	return body.DeviceID, body.Door, body.Delay, nil
}

func (rq *Request) DeviceDoorControl() (*uint32, *uint8, *string, error) {
	body := struct {
		DeviceID     *uint32 `json:"device-id"`
		Door         *uint8  `json:"door"`
		ControlState *string `json:"control"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, nil, errors.New("Missing device ID")
	}

	if body.Door == nil {
		return nil, nil, nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, nil, nil, errors.New("Invalid door")
	}

	if body.ControlState == nil {
		return nil, nil, nil, errors.New("Invalid door control state")
	} else if *body.ControlState != "normally open" && *body.ControlState != "normally closed" && *body.ControlState != "controlled" {
		return nil, nil, nil, errors.New("Invalid door control state")
	}

	return body.DeviceID, body.Door, body.ControlState, nil
}

func (rq *Request) DeviceCardID() (*uint32, *uint32, error) {
	body := struct {
		DeviceID   *uint32 `json:"device-id"`
		CardNumber *uint32 `json:"card-number"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, errors.New("Missing device ID")
	}

	if body.CardNumber == nil {
		return nil, nil, errors.New("Invalid card number")
	}

	return body.DeviceID, body.CardNumber, nil
}

func (rq *Request) DeviceCard() (*uint32, *types.Card, error) {
	body := struct {
		DeviceID   *uint32     `json:"device-id"`
		CardNumber *uint32     `json:"card-number"`
		From       *types.Date `json:"from"`
		To         *types.Date `json:"to"`
		Doors      []uint8     `json:"doors"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, errors.New("Missing device ID")
	}

	if body.CardNumber == nil {
		return nil, nil, errors.New("Invalid card number")
	}

	if body.From == nil {
		return nil, nil, errors.New("Invalid card 'from' date")
	}

	if body.To == nil {
		return nil, nil, errors.New("Invalid card 'to' date")
	}

	doors := []bool{false, false, false, false}
	for _, door := range body.Doors {
		if door >= 1 && door <= 4 {
			doors[door-1] = true
		}
	}

	return body.DeviceID, &types.Card{
		CardNumber: *body.CardNumber,
		From:       *body.From,
		To:         *body.To,
		Doors:      doors,
	}, nil
}

func (rq *Request) DeviceEventID() (*uint32, *uint32, error) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		EventID  *uint32 `json:"event-id"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	}

	if body.DeviceID == nil {
		return nil, nil, errors.New("Missing device ID")
	} else if *body.DeviceID == 0 {
		return nil, nil, errors.New("Missing device ID")
	}

	if body.EventID == nil {
		return nil, nil, errors.New("Invalid event ID")
	}

	return body.DeviceID, body.EventID, nil
}

func (rq *Request) DateRange() (*types.DateTime, *types.DateTime, error) {
	body := struct {
		From *startdate `json:"start"`
		To   *enddate   `json:"end"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	}

	if body.From != nil && body.To != nil {
		if time.Time(*body.To).Before(time.Time(*body.From)) {
			return nil, nil, fmt.Errorf("Invalid date range (%s to %s)", (*types.DateTime)(body.From), (*types.DateTime)(body.To))
		}
	}

	return (*types.DateTime)(body.From), (*types.DateTime)(body.To), nil
}

type startdate time.Time
type enddate time.Time

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
