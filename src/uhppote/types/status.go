package types

import (
	"fmt"
	"strings"
)

type Status struct {
	SerialNumber   SerialNumber
	LastIndex      uint32
	EventType      byte
	Granted        bool
	Door           byte
	DoorOpened     bool
	UserID         uint32
	EventTimestamp DateTime
	EventResult    byte
	DoorState      []bool
	DoorButton     []bool
	SystemState    byte
	SystemDateTime DateTime
	PacketNumber   uint32
	Backup         uint32
	SpecialMessage byte
	Battery        byte
	FireAlarm      byte
}

func (s *Status) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%s", s.SerialNumber))
	b.WriteString(fmt.Sprintf(" %-5d", s.LastIndex))
	b.WriteString(fmt.Sprintf(" %-3d", s.EventType))
	b.WriteString(fmt.Sprintf(" %-5v", s.Granted))
	b.WriteString(fmt.Sprintf(" %d", s.Door))
	b.WriteString(fmt.Sprintf(" %-5v", s.DoorOpened))
	b.WriteString(fmt.Sprintf(" %-10d", s.UserID))
	b.WriteString(fmt.Sprintf(" %s", s.EventTimestamp.String()))
	b.WriteString(fmt.Sprintf(" %-3d", s.EventResult))
	b.WriteString(fmt.Sprintf(" %-5v %-5v %-5v %-5v", s.DoorState[0], s.DoorState[1], s.DoorState[2], s.DoorState[3]))
	b.WriteString(fmt.Sprintf(" %-5v %-5v %-5v %-5v", s.DoorButton[0], s.DoorButton[1], s.DoorButton[2], s.DoorButton[3]))
	b.WriteString(fmt.Sprintf(" %d", s.SystemState))
	b.WriteString(fmt.Sprintf(" %s", s.SystemDateTime.String()))
	b.WriteString(fmt.Sprintf(" %d", s.PacketNumber))
	b.WriteString(fmt.Sprintf(" %d", s.Backup))
	b.WriteString(fmt.Sprintf(" %d", s.SpecialMessage))
	b.WriteString(fmt.Sprintf(" %v", s.Battery))
	b.WriteString(fmt.Sprintf(" %v", s.FireAlarm))

	return b.String()
}
