package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type DeviceTime struct {
	Device struct {
		ID       uint32         `json:"id"`
		DateTime types.DateTime `json:"date-time"`
	} `json:"device"`
}

type GetTimeRequest struct {
	DeviceID uint32
}

type GetTimeResponse struct {
	DeviceID uint32         `json:"device-id"`
	DateTime types.DateTime `json:"date-time"`
}

func (u *UHPPOTED) GetTime(ctx context.Context, request GetTimeRequest) (*GetTimeResponse, int, error) {
	u.debug(ctx, "get-time", request)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetTime(request.DeviceID)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetTimeResponse{
		DeviceID: uint32(result.SerialNumber),
		DateTime: result.DateTime,
	}

	u.debug(ctx, "get-time", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

func (u *UHPPOTED) SetTime(ctx context.Context, rq Request) {
	u.debug(ctx, "set-time", rq)

	id, err := rq.DeviceID()
	if err != nil {
		u.warn(ctx, 0, "set-time", err)
		u.oops(ctx, "set-time", "Error setting device time (invalid device ID)", StatusBadRequest)
		return
	}

	datetime, err := rq.DateTime()
	if err != nil || datetime == nil {
		u.warn(ctx, *id, "set-time", err)
		u.oops(ctx, "set-time", "Error setting device time (invalid date-time)", StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetTime(*id, *datetime)
	if err != nil {
		u.warn(ctx, *id, "set-time", err)
		u.oops(ctx, "set-time", "Error setting device time", StatusInternalServerError)
		return
	}

	response := DeviceTime{
		struct {
			ID       uint32         `json:"id"`
			DateTime types.DateTime `json:"date-time"`
		}{
			ID:       *id,
			DateTime: result.DateTime,
		},
	}

	u.reply(ctx, response)
}
