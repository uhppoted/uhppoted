package mqtt

import (
	"context"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"uhppoted"
)

func (m *MQTTD) getDevices(impl *uhppoted.UHPPOTED, ctx context.Context, msg MQTT.Message) {
	rq := uhppoted.GetDevicesRequest{}

	response, status, err := impl.GetDevices(ctx, rq)
	if err != nil {
		m.OnError(ctx, "get-devices", "Error retrieving list of devices", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetDevicesResponse
		}{
			MetaInfo:           getMetaInfo(ctx),
			GetDevicesResponse: *response,
		}

		m.Reply(ctx, reply)
	}
}
