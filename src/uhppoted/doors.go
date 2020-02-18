package uhppoted

import (
	"encoding/json"
	"fmt"
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
	DeviceID DeviceID
	Door     uint8
}

type GetDoorDelayResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Door     uint8    `json:"door"`
	Delay    uint8    `json:"delay"`
}

func (u *UHPPOTED) GetDoorDelay(request GetDoorDelayRequest) (*GetDoorDelayResponse, error) {
	u.debug("get-door-delay", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	result, err := u.Uhppote.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %v delay for %v (%w)", door, device, err))
	}

	response := GetDoorDelayResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug("get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type SetDoorDelayRequest struct {
	DeviceID DeviceID
	Door     uint8
	Delay    uint8
}

type SetDoorDelayResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Door     uint8    `json:"door"`
	Delay    uint8    `json:"delay"`
}

func (u *UHPPOTED) SetDoorDelay(request SetDoorDelayRequest) (*SetDoorDelayResponse, error) {
	u.debug("set-door-delay", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	state, err := u.Uhppote.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %s delay for %v (%w)", door, device, err))
	}

	result, err := u.Uhppote.SetDoorControlState(uint32(request.DeviceID), request.Door, state.ControlState, request.Delay)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error setting door %v delay %v for %v (%w)", door, state.ControlState, device, err))
	}

	response := SetDoorDelayResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Delay:    result.Delay,
	}

	u.debug("get-door-delay", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type GetDoorControlRequest struct {
	DeviceID DeviceID
	Door     uint8
}

type GetDoorControlResponse struct {
	DeviceID DeviceID     `json:"device-id"`
	Door     uint8        `json:"door"`
	Control  ControlState `json:"control"`
}

func (u *UHPPOTED) GetDoorControl(request GetDoorControlRequest) (*GetDoorControlResponse, error) {
	u.debug("get-door-control", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	result, err := u.Uhppote.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %s control for %v (%w)", door, device, err))
	}

	response := GetDoorControlResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Control:  ControlState(result.ControlState),
	}

	u.debug("get-door-control", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type SetDoorControlRequest struct {
	DeviceID DeviceID
	Door     uint8
	Control  ControlState
}

type SetDoorControlResponse struct {
	DeviceID DeviceID     `json:"device-id"`
	Door     uint8        `json:"door"`
	Control  ControlState `json:"control"`
}

func (u *UHPPOTED) SetDoorControl(request SetDoorControlRequest) (*SetDoorControlResponse, error) {
	u.debug("set-door-control", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	door := request.Door
	state, err := u.Uhppote.GetDoorControlState(device, door)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting door %s control for %v (%w)", door, device, err))
	}

	result, err := u.Uhppote.SetDoorControlState(device, door, uint8(request.Control), state.Delay)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error setting door %s control %v for %v (%w)", door, request.Control, device, err))
	}

	response := SetDoorControlResponse{
		DeviceID: DeviceID(result.SerialNumber),
		Door:     result.Door,
		Control:  ControlState(result.ControlState),
	}

	u.debug("set-door-control", fmt.Sprintf("response %+v", response))

	return &response, nil
}
