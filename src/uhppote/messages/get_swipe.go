package messages

import (
	"encoding/binary"
	"fmt"
	"uhppote/types"
)

type GetSwipe struct {
	StartOfMessage byte
	MsgType        byte
	SerialNumber   uint32
	Swipe          *types.Swipe
}

func NewGetSwipe(msg []byte) (*GetSwipe, error) {
	index := binary.LittleEndian.Uint32(msg[8:12])
	var swipe *types.Swipe

	if index > 0 {
		timestamp, err := types.DecodeDateTime(msg[20:27])
		if err != nil {
			panic(fmt.Sprintf("Unexpected error decoding timestamp: [%v]", err))
		}

		swipe = &types.Swipe{
			Index:      binary.LittleEndian.Uint32(msg[8:12]),
			Type:       msg[12],
			Access:     msg[13] == 0x01,
			Door:       msg[14],
			DoorState:  msg[15],
			Card:       binary.LittleEndian.Uint32(msg[16:20]),
			Timestamp:  *timestamp,
			RecordType: msg[27],
		}
	}

	return &GetSwipe{
		StartOfMessage: msg[0],
		MsgType:        msg[1],
		SerialNumber:   binary.LittleEndian.Uint32(msg[4:8]),
		Swipe:          swipe,
	}, nil
}
