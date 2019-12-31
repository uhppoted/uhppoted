package mqtt

import (
	"encoding/json"
	"errors"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
