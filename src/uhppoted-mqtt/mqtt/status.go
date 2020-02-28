package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppoted-api/uhppoted"
)

func (m *MQTTD) getStatus(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	rq := uhppoted.GetStatusRequest{
		DeviceID: *body.DeviceID,
	}

	response, err := impl.GetStatus(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving status for %v", *body.DeviceID))
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
