package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDevices(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, status, err := impl.GetDevices(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving list of devices", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetDevicesResponse
	}{
		metainfo:           meta,
		GetDevicesResponse: *response,
	}, nil
}

func (m *MQTTD) getDevice(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil, nil
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetDevice(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	reply := struct {
		metainfo
		uhppoted.GetDeviceResponse
	}{
		metainfo:          meta,
		GetDeviceResponse: *response,
	}

	return reply, nil
}
