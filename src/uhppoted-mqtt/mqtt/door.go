package mqtt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDoorDelay(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
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

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door",
		}
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, status, err := impl.GetDoorDelay(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error retrieving door %v delay", *body.Door),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetDoorDelayResponse
	}{
		metainfo:             meta,
		GetDoorDelayResponse: *response,
	}, nil
}

func (m *MQTTD) setDoorDelay(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
		Delay    *uint8             `json:"delay"`
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

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door",
		}
	}

	if body.Delay == nil || *body.Delay == 0 || *body.Delay > 60 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door delay"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door delay",
		}
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Delay:    *body.Delay,
	}

	response, status, err := impl.SetDoorDelay(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error setting door %v delay", *body.Door),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.SetDoorDelayResponse
	}{
		metainfo:             meta,
		SetDoorDelayResponse: *response,
	}, nil
}

func (m *MQTTD) getDoorControl(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
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

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door",
		}
	}

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, status, err := impl.GetDoorControl(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error getting door %v control", *body.Door),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.GetDoorControlResponse
	}{
		metainfo:               meta,
		GetDoorControlResponse: *response,
	}, nil
}

func (m *MQTTD) setDoorControl(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID     `json:"device-id"`
		Door     *uint8                 `json:"door"`
		Control  *uhppoted.ControlState `json:"control"`
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

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door",
		}
	}

	if body.Control == nil || *body.Control < 1 || *body.Control > 3 {
		return nil, &errorx{
			Err:     errors.New("Missing/invalid door control value"),
			Code:    uhppoted.StatusBadRequest,
			Message: "Missing/invalid door control value",
		}
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Control:  *body.Control,
	}

	response, status, err := impl.SetDoorControl(ctx, rq)
	if err != nil {
		return nil, &errorx{
			Err:     err,
			Code:    status,
			Message: fmt.Sprintf("Error setting door %v control", *body.Door),
		}
	}

	if response == nil {
		return nil, nil
	}

	return struct {
		metainfo
		uhppoted.SetDoorControlResponse
	}{
		metainfo:               meta,
		SetDoorControlResponse: *response,
	}, nil
}
