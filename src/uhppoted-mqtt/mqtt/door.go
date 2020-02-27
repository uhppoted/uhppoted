package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppoted"
)

func (m *MQTTD) getDoorDelay(meta metainfo, impl *uhppoted.UHPPOTED, ctx context.Context, request []byte) (interface{}, error) {
	body := struct {
		DeviceID *uhppoted.DeviceID `json:"device-id"`
		Door     *uint8             `json:"door"`
	}{}

	if err := json.Unmarshal(request, &body); err != nil {
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, InvalidDoorID
	}

	rq := uhppoted.GetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorDelay(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error retrieving door %v delay", *body.Door))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, InvalidDoorID
	}

	if body.Delay == nil || *body.Delay == 0 || *body.Delay > 60 {
		return nil, InvalidDoorDelay
	}

	rq := uhppoted.SetDoorDelayRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Delay:    *body.Delay,
	}

	response, err := impl.SetDoorDelay(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error setting door %v delay", *body.Door))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, InvalidDoorID
	}

	rq := uhppoted.GetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
	}

	response, err := impl.GetDoorControl(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error getting door %v control", *body.Door))
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
		return nil, ferror(fmt.Errorf("%w: %v", uhppoted.BadRequest, err), "Cannot parse request")
	}

	if body.DeviceID == nil {
		return nil, InvalidDeviceID
	}

	if body.Door == nil || *body.Door < 1 || *body.Door > 4 {
		return nil, InvalidDoorID
	}

	if body.Control == nil || *body.Control < 1 || *body.Control > 3 {
		return nil, InvalidDoorControl
	}

	rq := uhppoted.SetDoorControlRequest{
		DeviceID: *body.DeviceID,
		Door:     *body.Door,
		Control:  *body.Control,
	}

	response, err := impl.SetDoorControl(rq)
	if err != nil {
		return nil, ferror(err, fmt.Sprintf("Error setting door %v control", *body.Door))
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
