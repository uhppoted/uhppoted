package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"uhppoted"
)

func (m *MQTTD) getStatus(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	body := struct {
		DeviceID *uint32 `json:"device-id"`
	}{}

	if err := json.Unmarshal(msg.Payload(), &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil || *body.DeviceID == 0 {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(msg.Payload())))
		return
	}

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetStatus(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device status", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetStatusResponse
		}{
			MetaInfo:          getMetaInfo(ctx),
			GetStatusResponse: *response,
		}

		m.reply(ctx, reply)
	}
}
