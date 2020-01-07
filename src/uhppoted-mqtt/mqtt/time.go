package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getTime(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) {
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

	rq := uhppoted.GetTimeRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetTime(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving current device time", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.GetTimeResponse
		}{
			MetaInfo:        getMetaInfo(ctx),
			GetTimeResponse: *response,
		}

		m.reply(ctx, reply)
	}
}

func (m *MQTTD) setTime(impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		DateTime *types.DateTime    `json:"date-time"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return
	}

	if body.DateTime == nil {
		m.OnError(ctx, "Missing/invalid device time", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device time '%s'", string(request)))
		return
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: *body.DeviceID,
		DateTime: *body.DateTime,
	}

	response, status, err := impl.SetTime(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error setting current device time", status, err)
		return
	}

	if response != nil {
		reply := struct {
			MetaInfo *metainfo `json:"meta-info,omitempty"`
			uhppoted.SetTimeResponse
		}{
			MetaInfo:        getMetaInfo(ctx),
			SetTimeResponse: *response,
		}

		m.reply(ctx, reply)
	}
}
