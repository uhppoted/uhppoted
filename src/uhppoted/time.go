package uhppoted

import (
	"fmt"
	"time"
	"uhppote/types"
)

type GetTimeRequest struct {
	DeviceID DeviceID
}

type GetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) GetTime(request GetTimeRequest) (*GetTimeResponse, error) {
	u.debug("get-time", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	result, err := u.Uhppote.GetTime(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting time for %v (%w)", device, err))
	}

	response := GetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("get-time", fmt.Sprintf("response %+v", response))

	return &response, nil
}

type SetTimeRequest struct {
	DeviceID DeviceID
	DateTime types.DateTime
}

type SetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) SetTime(request SetTimeRequest) (*SetTimeResponse, error) {
	u.debug("set-time", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	result, err := u.Uhppote.SetTime(device, time.Time(request.DateTime))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error setting time for %v (%w)", device, err))
	}

	response := SetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("set-time", fmt.Sprintf("response %+v", response))

	return &response, nil
}
