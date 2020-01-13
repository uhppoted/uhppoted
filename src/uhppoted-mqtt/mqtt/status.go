package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getStatus(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil
	}

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetStatus(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device status", status, err)
		return nil
	}

	if response == nil {
		return nil
	}

	return struct {
		metainfo
		uhppoted.GetStatusResponse
	}{
		metainfo:          meta,
		GetStatusResponse: *response,
	}
}
