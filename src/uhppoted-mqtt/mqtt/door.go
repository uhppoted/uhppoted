package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"uhppoted"
)

func (m *MQTTD) getDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		Door     *uint8  `json:"door"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, status, err := impl.GetDoorDelay(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device door delay", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetDoorDelayResponse
		}{
			MetaInfo:             getMetaInfo(ctx),
			GetDoorDelayResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}

func (m *MQTTD) setDoorDelay(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
		Door     *uint8  `json:"door"`
		Delay    *uint8  `json:"delay"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(msg.Payload())))
		return
	}

	if body.Delay == nil || *body.Delay == 0 || *body.Delay > 60 {
		m.OnError(ctx, "Missing/invalid device door delay value", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door delay value '%s'", string(msg.Payload())))
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Delay:    *body.Delay,
	}

	response, status, err := impl.SetDoorDelay(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error setting device door delay", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.SetDoorDelayResponse
		}{
			MetaInfo:             getMetaInfo(ctx),
			SetDoorDelayResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}
