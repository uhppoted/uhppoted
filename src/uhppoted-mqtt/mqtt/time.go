package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"uhppote/types"
	"uhppoted"
)

func (m *MQTTD) getTime(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
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

	rq := uhppoted.GetTimeRequest{
		DeviceID: *body.DeviceID,
	}

	response, status, err := impl.GetTime(rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: "Error retrieving device time",
		}
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

	if body.DateTime == nil {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid date-time"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid date-time",
		}
	}

	rq := uhppoted.SetTimeRequest{
		DeviceID: *body.DeviceID,
		DateTime: *body.DateTime,
	}

	response, status, err := impl.SetTime(rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: "Error setting device date/time",
		}
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
