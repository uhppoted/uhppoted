package uhppoted

import (
	"context"
	"fmt"
	"time"
	"uhppote"
	"uhppote/types"
)

type GetTimeRequest struct {
	DeviceID DeviceID
}

type GetTimeResponse struct {
	DeviceID DeviceID       `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) GetTime(ctx context.Context, request GetTimeRequest) (*GetTimeResponse, int, error) {
	u.debug("get-time", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetTime(uint32(request.DeviceID))
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
	DeviceID uint32
	DateTime types.DateTime
}

type SetTimeResponse struct {
	DeviceID uint32         `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) SetTime(ctx context.Context, request SetTimeRequest) (*SetTimeResponse, int, error) {
	u.debug("set-time", fmt.Sprintf("request  %+v", request))

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetTime(request.DeviceID, time.Time(request.DateTime))
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := SetTimeResponse{
		DeviceID: uint32(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug("set-time", fmt.Sprintf("response %+v", response))

	return &response, StatusOK, nil
}
