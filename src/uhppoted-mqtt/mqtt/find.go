package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDevices(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, err := impl.GetDevices(rq)
	if err != nil {
		return nil, ferror(err, "Error searching for active devices")
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetDevice(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Could not retrieve device information for %d", *body.DeviceID))
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
