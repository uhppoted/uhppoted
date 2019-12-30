package uhppoted

import (
	"context"
	"fmt"
	"uhppote"
	"uhppote/types"
)

type Status struct {
	LastEventIndex uint32         `json:"last-event-index"`
	EventType      byte           `json:"event-type"`
	Granted        bool           `json:"access-granted"`
	Door           byte           `json:"door"`
	DoorOpened     bool           `json:"door-opened"`
	UserID         uint32         `json:"user-id"`
	EventTimestamp types.DateTime `json:"event-timestamp"`
	EventResult    byte           `json:"event-result"`
	DoorState      []bool         `json:"door-states"`
	DoorButton     []bool         `json:"door-buttons"`
	SystemState    byte           `json:"system-state"`
	SystemDateTime types.DateTime `json:"system-datetime"`
	PacketNumber   uint32         `json:"packet-number"`
	Backup         uint32         `json:"backup-state"`
	SpecialMessage byte           `json:"special-message"`
	Battery        byte           `json:"battery-status"`
	FireAlarm      byte           `json:"fire-alarm-status"`
}

type GetStatusRequest struct {
	DeviceID uint32
}

type GetStatusResponse struct {
	DeviceID uint32 `json:"device-id"`
	Status   Status `json:"status"`
}

func (u *UHPPOTED) GetStatus(ctx context.Context, request GetStatusRequest) (*GetStatusResponse, int, error) {
	u.debug(ctx, "get-status", request)

	status, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetStatus(request.DeviceID)
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	response := GetStatusResponse{
		DeviceID: uint32(status.SerialNumber),
		Status: Status{
			LastEventIndex: status.LastIndex,
			EventType:      status.EventType,
			Granted:        status.Granted,
			Door:           status.Door,
			DoorOpened:     status.DoorOpened,
			UserID:         status.UserID,
			EventTimestamp: status.EventTimestamp,
			EventResult:    status.EventResult,
			DoorState:      status.DoorState,
			DoorButton:     status.DoorButton,
			SystemState:    status.SystemState,
			SystemDateTime: status.SystemDateTime,
			PacketNumber:   status.PacketNumber,
			Backup:         status.Backup,
			SpecialMessage: status.SpecialMessage,
			Battery:        status.Battery,
			FireAlarm:      status.FireAlarm,
		},
	}

	u.debug(ctx, "get-status", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}
