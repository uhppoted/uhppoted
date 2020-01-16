package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getStatus(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
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

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetStatus(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    uhppoted.StatusInternalServerError,
			Message: fmt.Sprintf("Could not retrieve status for %d", *body.DeviceID),
		}
	}

	if response != nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetStatusResponse
	}{
		metainfo:          meta,
		GetStatusResponse: *response,
	}, nil
}
