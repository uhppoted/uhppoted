package types

import (
	"fmt"
	"strings"
)

type Status struct {
	SerialNumber   SerialNumber
	LastIndex      uint32
	SwipeRecord    byte
	Granted        bool
	Door           byte
	DoorOpened     bool
	UserId         uint32
	SwipeDateTime  DateTime
	SwipeReason    byte
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
	b.WriteString(fmt.Sprintf(" %d", s.LastIndex))
	b.WriteString(fmt.Sprintf(" %d", s.SwipeRecord))
	b.WriteString(fmt.Sprintf(" %v", s.Granted))
	b.WriteString(fmt.Sprintf(" %d", s.Door))
	b.WriteString(fmt.Sprintf(" %v", s.DoorOpened))
	b.WriteString(fmt.Sprintf(" %d", s.UserId))
	b.WriteString(fmt.Sprintf(" %s", s.SwipeDateTime.String()))
	b.WriteString(fmt.Sprintf(" %d", s.SwipeReason))
	b.WriteString(fmt.Sprintf(" %v %v %v %v", s.DoorState[0], s.DoorState[1], s.DoorState[2], s.DoorState[3]))
	b.WriteString(fmt.Sprintf(" %v %v %v %v", s.DoorButton[0], s.DoorButton[1], s.DoorButton[2], s.DoorButton[3]))
	b.WriteString(fmt.Sprintf(" %d", s.SystemState))
	b.WriteString(fmt.Sprintf(" %s", s.SystemDateTime.String()))
	b.WriteString(fmt.Sprintf(" %d", s.PacketNumber))
	b.WriteString(fmt.Sprintf(" %d", s.Backup))
	b.WriteString(fmt.Sprintf(" %d", s.SpecialMessage))
	b.WriteString(fmt.Sprintf(" %v", s.Battery))
	b.WriteString(fmt.Sprintf(" %v", s.FireAlarm))

	return b.String()
}

//func DecodeStatus(bytes []byte) (*Status, error) {
//	serialNumber := binary.LittleEndian.Uint32(bytes[4:8])
//	lastIndex := binary.LittleEndian.Uint32(bytes[8:12])
//	swipeRecord := bytes[12]
//	granted := bytes[13] == 0x01
//	door := bytes[14]
//	doorOpen := bytes[15] == 0x01
//	cardNumber := binary.LittleEndian.Uint32(bytes[16:20])
//
//	datetime, err := DecodeDateTime(bytes[20:27])
//	if err != nil {
//		return nil, err
//	}
//
//	reason := bytes[27]
//
//	door1 := bytes[28] == 0x01
//	door2 := bytes[29] == 0x01
//	door3 := bytes[30] == 0x01
//	door4 := bytes[31] == 0x01
//
//	button1 := bytes[32] == 0x01
//	button2 := bytes[33] == 0x01
//	button3 := bytes[34] == 0x01
//	button4 := bytes[35] == 0x01
//
//	systemState := bytes[36]
//	systemTime, err := DecodeDateTime([]byte{0x20, bytes[51], bytes[52], bytes[53], bytes[37], bytes[38], bytes[39]})
//	if err != nil {
//		return nil, err
//	}
//
//	packetNumber := binary.LittleEndian.Uint32(bytes[40:44])
//	backup := binary.LittleEndian.Uint32(bytes[44:48])
//	specialMessage := bytes[48]
//	lowBattery := bytes[49] == 0x01
//	fireAlarm := bytes[50] == 0x01
//
//	return &Status{
//		SerialNumber:   SerialNumber(serialNumber),
//		LastIndex:      lastIndex,
//		SwipeRecord:    swipeRecord,
//		Granted:        granted,
//		Door:           door,
//		DoorOpen:       doorOpen,
//		CardNumber:     cardNumber,
//		SwipeDateTime:  *datetime,
//		SwipeReason:    reason,
//		DoorState:      []bool{door1, door2, door3, door4},
//		DoorButton:     []bool{button1, button2, button3, button4},
//		SystemState:    systemState,
//		SystemDateTime: *systemTime,
//		PacketNumber:   packetNumber,
//		Backup:         backup,
//		SpecialMessage: specialMessage,
//		LowBattery:     lowBattery,
//		FireAlarm:      fireAlarm,
//	}, nil
//}
