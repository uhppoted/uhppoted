package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getStatus(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) {
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
