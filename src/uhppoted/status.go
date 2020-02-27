package uhppoted

import (
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
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
	DeviceID DeviceID
}

type GetStatusResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Status   Status   `json:"status"`
}

func (u *UHPPOTED) GetStatus(request GetStatusRequest) (*GetStatusResponse, error) {
	u.debug("get-status", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	status, err := u.Uhppote.GetStatus(device)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error retrieving status for %v (%w)", device, err))
	}

	response := GetStatusResponse{
		DeviceID: DeviceID(status.SerialNumber),
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

	u.debug("get-status", fmt.Sprintf("response %+v", response))

	return &response, nil
}
