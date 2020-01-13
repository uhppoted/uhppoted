package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppoted"
)

func (m *MQTTD) getDoorDelay(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(request)))
		return nil
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, status, err := impl.GetDoorDelay(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device door delay", status, err)
		return nil
	}

	if response == nil {
		return nil
	}

	return struct {
		metainfo
		uhppoted.GetDoorDelayResponse
	}{
		metainfo:             meta,
		GetDoorDelayResponse: *response,
	}
}

func (m *MQTTD) setDoorDelay(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
		Delay    *uint8             `json:"delay"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing device ID: %s", string(request)))
		return nil
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(request)))
		return nil
	}

	if body.Delay == nil || *body.Delay == 0 || *body.Delay > 60 {
		m.OnError(ctx, "Missing/invalid device door delay value", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door delay value '%s'", string(request)))
		return nil
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Delay:    *body.Delay,
	}

	response, status, err := impl.SetDoorDelay(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error setting device door delay", status, err)
		return nil
	}

	if response == nil {
		return nil
	}

	return struct {
		metainfo
		uhppoted.SetDoorDelayResponse
	}{
		metainfo:             meta,
		SetDoorDelayResponse: *response,
	}
}

func (m *MQTTD) getDoorControl(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(request)))
		return nil
	}

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, status, err := impl.GetDoorControl(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error retrieving device door control", status, err)
		return nil
	}

	if response == nil {
		return nil
	}

	return struct {
		metainfo
		uhppoted.GetDoorControlResponse
	}{
		metainfo:               meta,
		GetDoorControlResponse: *response,
	}
}

func (m *MQTTD) setDoorControl(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) interface{} {
	body := struct {
		DeviceID *uhppoted.DeviceID     `json:"device-id"`
		Door     *uint8                 `json:"door"`
		Control  *uhppoted.ControlState `json:"control"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		m.OnError(ctx, "Cannot parse request", uhppoted.StatusBadRequest, err)
		return nil
	}

	if body.DeviceID == nil {
		m.OnError(ctx, "Missing/invalid device ID", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device ID '%s'", string(request)))
		return nil
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		m.OnError(ctx, "Missing/invalid device door", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door '%s'", string(request)))
		return nil
	}

	if body.Control == nil || *body.Control < 1 || *body.Control > 3 {
		m.OnError(ctx, "Missing/invalid device door control value", uhppoted.StatusBadRequest, fmt.Errorf("Missing/invalid device door control value '%s'", string(request)))
		return nil
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Control:  *body.Control,
	}

	response, status, err := impl.SetDoorControl(ctx, rq)
	if err != nil {
		m.OnError(ctx, "Error setting device door control", status, err)
		return nil
	}

	if response == nil {
		return nil
	}
	return struct {
		metainfo
		uhppoted.SetDoorControlResponse
	}{
		metainfo:               meta,
		SetDoorControlResponse: *response,
	}
}
