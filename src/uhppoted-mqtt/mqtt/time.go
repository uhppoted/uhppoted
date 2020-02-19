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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	rq := uhppoted.GetTimeRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetTime(rq)
	if err != nil {
		return nil, ferror(err, "Error retrieving device time")
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.DateTime == nil {
		return nil, InvalidDateTime
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: *body.DeviceID,
		DateTime: *body.DateTime,
	}

	response, err := impl.SetTime(rq)
	if err != nil {
		return nil, ferror(err, "Error setting device date/time")
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
