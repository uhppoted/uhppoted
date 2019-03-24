package messages

import (
	"encoding/binary"
)

type OpenDoorReply struct {
	StartOfMessage byte
	MsgType        byte
	SerialNumber   uint32
	Opened         bool
}

func NewOpenDoorReply(msg []byte) (*OpenDoorReply, error) {
	serialNumber := binary.LittleEndian.Uint32(msg[4:8])
	opened := msg[8] == 0x01

	return &OpenDoorReply{
		msg[0],
		msg[1],
		serialNumber,
		opened,
	}, nil
}
