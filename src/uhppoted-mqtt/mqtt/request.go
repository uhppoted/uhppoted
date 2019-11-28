package mqtt

import (
	"encoding/json"
	"errors"
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

func (rq *Request) DeviceId() (*uint32, error) {
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

func (rq *Request) DoorDelay() (*uint8, *uint8, error) {
	body := struct {
		Door  *uint8 `json:"door"`
		Delay *uint8 `json:"delay"`
	}{}

	if err := json.Unmarshal(rq.Message.Payload(), &body); err != nil {
		return nil, nil, err
	} else if body.Door == nil {
		return nil, nil, errors.New("Invalid door")
	} else if *body.Door < 1 || *body.Door > 4 {
		return nil, nil, errors.New("Invalid door")
	} else if body.Delay == nil {
		return nil, nil, errors.New("Invalid door delay")
	} else if *body.Delay == 0 || *body.Delay > 60 {
		return nil, nil, errors.New("Invalid door delay")
	}

	return body.Door, body.Delay, nil
}
