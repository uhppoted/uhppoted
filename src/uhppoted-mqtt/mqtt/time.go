package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getTime(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
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

	rq := uhppoted.GetTimeRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetTime(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving current device time", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetTimeResponse
	}{
		metainfo:        meta,
		GetTimeResponse: *response,
	}, nil
}

func (m *MQTTD) setTime(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		DateTime *types.DateTime    `json:"date-time"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil, nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil, nil
	}

	if body.DateTime == nil {
		m.OnError(ctx, "Missing/invalid device time", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device time '%s'", string(request)))
		return nil, nil
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: *body.DeviceID,
		DateTime: *body.DateTime,
	}

	response, status, err := impl.SetTime(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error setting current device time", status, err)
		return nil, nil
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.SetTimeResponse
	}{
		metainfo:        meta,
		SetTimeResponse: *response,
	}, nil
}
