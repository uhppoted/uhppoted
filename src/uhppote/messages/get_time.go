package messages

import (
	"encoding/binary"
	"uhppote/types"
)

type GetTime struct {
	StartOfMessage byte
	MsgType        byte
	SerialNumber   uint32
	DateTime       types.DateTime
}

func NewGetTime(msg []byte) (*GetTime, error) {
	datetime, err := types.DecodeDateTime(msg[8:15])
	if err != nil {
		return nil, err
	}

	return &GetTime{
		StartOfMessage: msg[0],
		MsgType:        msg[1],
		SerialNumber:   binary.LittleEndian.Uint32(msg[4:8]),
		DateTime:       *datetime,
	}, nil
}
