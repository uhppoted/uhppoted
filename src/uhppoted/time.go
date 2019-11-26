package uhppoted

import (
	"context"
	"uhppote"
	"uhppote/types"
)

type DeviceTime struct {
	Device struct {
		ID       uint32         `json:"id"`
		DateTime types.DateTime `json:"date-time"`
	} `json:"device"`
}

func (u *UHPPOTED) GetTime(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-time", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, id, "get-time", err)
		u.oops(ctx, "get-time", "Error retrieving device time (invalid device ID)", StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetTime(id)
	if err != nil {
		u.warn(ctx, id, "get-time", err)
		u.oops(ctx, "get-status", "Error retrieving device time", StatusInternalServerError)
		return
	}

	response := DeviceTime{
		struct {
			ID       uint32         `json:"id"`
			DateTime types.DateTime `json:"date-time"`
		}{
			ID:       id,
			DateTime: result.DateTime,
		},
	}

	u.reply(ctx, response)
}
