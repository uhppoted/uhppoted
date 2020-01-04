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
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	request, ok := ctx.Value("request").([]byte)
	if !ok {
		m.OnError(ctx, "Message payload does not contain 'request' field", uhppoted.StatusBadRequest, fmt.Errorf("Bad Request"))
		return
	}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil {
		m.OnError(ctx, string(msg.Payload()), uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID"))
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
