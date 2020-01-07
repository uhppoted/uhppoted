package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDevices(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) {
	rq := uhppoted.GetDevicesRequest{}

	response, status, err := impl.GetDevices(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving list of devices", status, err)
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

		m.reply(ctx, reply)
	}
}

func (m *MQTTD) getDevice(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetDevice(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetDeviceResponse
		}{
			MetaInfo:          getMetaInfo(ctx),
			GetDeviceResponse: *response,
		}

		m.reply(ctx, reply)
	}
}
