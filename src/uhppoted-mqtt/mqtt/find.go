package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDevices(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	rq := uhppoted.GetDevicesRequest{}

	response, status, err := impl.GetDevices(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: "Error searching for active devices",
		}
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
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusBadRequest,
			Message: "Cannot parse request",
		}
	}

	if body.DeviceID == nil {
		return nil, &errorx{
			Err:     errors.New("Missing device ID"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing device ID",
		}
	}

	rq := uhppoted.GetDeviceRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetDevice(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Could not retrieve device information for %d", *body.DeviceID),
		}
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
