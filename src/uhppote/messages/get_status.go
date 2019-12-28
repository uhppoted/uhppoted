package messages

import (
	"uhppote/types"
)

type GetStatusRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x20"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetStatusResponse struct {
	MsgType        types.MsgType      `uhppote:"value:0x20"`
	SerialNumber   types.SerialNumber `uhppote:"offset:4"`
	LastIndex      uint32             `uhppote:"offset:8"`
	EventType      byte               `uhppote:"offset:12"`
	Granted        bool               `uhppote:"offset:13"`
	Door           byte               `uhppote:"offset:14"`
	DoorOpened     bool               `uhppote:"offset:15"`
	UserID         uint32             `uhppote:"offset:16"`
	EventTimestamp types.DateTime     `uhppote:"offset:20"`
	EventResult    byte               `uhppote:"offset:27"`
	Door1State     bool               `uhppote:"offset:28"`
	Door2State     bool               `uhppote:"offset:29"`
	Door3State     bool               `uhppote:"offset:30"`
	Door4State     bool               `uhppote:"offset:31"`
	Door1Button    bool               `uhppote:"offset:32"`
	Door2Button    bool               `uhppote:"offset:33"`
	Door3Button    bool               `uhppote:"offset:34"`
	Door4Button    bool               `uhppote:"offset:35"`
	SystemState    byte               `uhppote:"offset:36"`
	SystemDate     types.SystemDate   `uhppote:"offset:51"`
	SystemTime     types.SystemTime   `uhppote:"offset:37"`
	PacketNumber   uint32             `uhppote:"offset:40"` // TODO unverified - trust at own risk
	Backup         uint32             `uhppote:"offset:44"` // TODO unverified - trust at own risk
	SpecialMessage byte               `uhppote:"offset:48"` // TODO unverified - trust at own risk
	Battery        byte               `uhppote:"offset:49"` // TODO unverified - trust at own risk
	FireAlarm      byte               `uhppote:"offset:50"` // TODO unverified - trust at own risk
}
