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

func (m *MQTTD) getDevice(operation string, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		RequestID *string            `json:"request-id"`
		DeviceID  *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetDevice(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device", status, err)
		return nil
	} else if response == nil {
		return nil
	}

	reply := struct {
		RequestID *string `json:"request-id,omitempty"`
		Operation string  `json:"operation,omitempty"`
		uhppoted.GetDeviceResponse
	}{
		RequestID:         body.RequestID,
		Operation:         operation,
		GetDeviceResponse: *response,
	}

	return reply
}
