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

func (u *UHPPOTED) GetTime(request GetTimeRequest) (*GetTimeResponse, int, error) {
	u.debug("get-time", fmt.Sprintf("request  %+v", request))

	result, err := u.Uhppote.GetTime(uint32(request.DeviceID))
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("get-time", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}

type SetTimeRequest struct {
	DeviceID DeviceID
	DateTime types.DateTime
}

type SetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) SetTime(request SetTimeRequest) (*SetTimeResponse, int, error) {
	u.debug("set-time", fmt.Sprintf("request  %+v", request))

	result, err := u.Uhppote.SetTime(uint32(request.DeviceID), time.Time(request.DateTime))
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := SetTimeResponse{
		DeviceID: DeviceID(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("set-time", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}
