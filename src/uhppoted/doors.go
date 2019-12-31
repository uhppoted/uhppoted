package uhppoted

import (
	"context"
	"encoding/json"
	"fmt"
	"uhppote"
)

type ControlState uint8

const (
	NormallyOpen   ControlState = 1
	NormallyClosed ControlState = 2
	Controlled     ControlState = 3
)

func (s ControlState) MarshalJSON() ([]byte, error) {
	switch s {
	case NormallyOpen:
		return json.Marshal("normally open")
	case NormallyClosed:
		return json.Marshal("normally closed")
	case Controlled:
		return json.Marshal("controlled")
	}

	return []byte("???"), fmt.Errorf("Invalid ControlState: %v", s)
}

func (s *ControlState) UnmarshalJSON(bytes []byte) (err error) {
	v := ""
	if err = json.Unmarshal(bytes, &v); err == nil {
		switch v {
		case "normally open":
			*s = NormallyOpen
		case "normally closed":
			*s = NormallyClosed
		case "controlled":
			*s = Controlled
		default:
			err = fmt.Errorf("Invalid DoorControlState: %s", string(bytes))
		}
	}

	return
}

type GetDoorDelayRequest struct {
	DeviceID uint32
	Door     uint8
}

type GetDoorDelayResponse struct {
	DeviceID uint32 `json:"device-id"`
	Door     uint8  `json:"door"`
	Delay    uint8  `json:"delay"`
}

func (u *UHPPOTED) GetDoorDelay(ctx context.Context, request GetDoorDelayRequest) (*GetDoorDelayResponse, int, error) {
	u.debug(ctx, "get-door-delay", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetDoorDelayResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug(ctx, "get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type SetDoorDelayRequest struct {
	DeviceID uint32
	Door     uint8
	Delay    uint8
}

type SetDoorDelayResponse struct {
	DeviceID uint32 `json:"device-id"`
	Door     uint8  `json:"door"`
	Delay    uint8  `json:"delay"`
}

func (u *UHPPOTED) SetDoorDelay(ctx context.Context, request SetDoorDelayRequest) (*SetDoorDelayResponse, int, error) {
	u.debug(ctx, "set-door-delay", fmt.Sprintf("request  %+v", request))

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(request.DeviceID, request.Door, state.ControlState, request.Delay)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := SetDoorDelayResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug(ctx, "get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type GetDoorControlRequest struct {
	DeviceID uint32
	Door     uint8
}

type GetDoorControlResponse struct {
	DeviceID uint32       `json:"device-id"`
	Door     uint8        `json:"door"`
	Control  ControlState `json:"control"`
}

func (u *UHPPOTED) GetDoorControl(ctx context.Context, request GetDoorControlRequest) (*GetDoorControlResponse, int, error) {
	u.debug(ctx, "get-door-control", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetDoorControlResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Control:  ControlState(result.ControlState),
	}

	u.debug(ctx, "get-door-control", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type SetDoorControlRequest struct {
	DeviceID uint32
	Door     uint8
	Control  ControlState
}

type SetDoorControlResponse struct {
	DeviceID uint32       `json:"device-id"`
	Door     uint8        `json:"door"`
	Control  ControlState `json:"control"`
}

func (u *UHPPOTED) SetDoorControl(ctx context.Context, request SetDoorControlRequest) (*SetDoorControlResponse, int, error) {
	u.debug(ctx, "set-door-control", fmt.Sprintf("request  %+v", request))

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(request.DeviceID, request.Door)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(request.DeviceID, request.Door, uint8(request.Control), state.Delay)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := SetDoorControlResponse{
		DeviceID: uint32(result.SerialNumber),
		Door:     result.Door,
		Control:  ControlState(result.ControlState),
	}

	u.debug(ctx, "set-door-control", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}
