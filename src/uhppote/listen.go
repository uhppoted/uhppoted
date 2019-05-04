package uhppote

import (
	"fmt"
	"time"
	"uhppote/types"
)

func (u *UHPPOTE) Listen() error {
	msg := GetStatusResponse{}

	err := u.listen(&msg)
	if err == nil {
		d := msg.SystemDate.Date.Format("2006-01-02")
		t := msg.SystemTime.Time.Format("15:04:05")
		datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", d+" "+t, time.Local)

		event := types.Status{
			SerialNumber:   msg.SerialNumber,
			LastIndex:      msg.LastIndex,
			SwipeRecord:    msg.SwipeRecord,
			Granted:        msg.Granted,
			Door:           msg.Door,
			DoorOpen:       msg.DoorOpen,
			CardNumber:     msg.CardNumber,
			SwipeDateTime:  msg.SwipeDateTime,
			SwipeReason:    msg.SwipeReason,
			DoorState:      []bool{msg.Door1State, msg.Door2State, msg.Door3State, msg.Door4State},
			DoorButton:     []bool{msg.Door1Button, msg.Door2Button, msg.Door3Button, msg.Door4Button},
			SystemState:    msg.SystemState,
			SystemDateTime: types.DateTime{DateTime: datetime},
			PacketNumber:   msg.PacketNumber,
			Backup:         msg.Backup,
			SpecialMessage: msg.SpecialMessage,
			LowBattery:     msg.LowBattery,
			FireAlarm:      msg.FireAlarm,
		}

		fmt.Printf("%s\n", event.String())
	}

	return err
}
