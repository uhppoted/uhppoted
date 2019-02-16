package messages

import (
	"encoding/binary"
	"uhppote/types"
)

type GetAuthRec struct {
	StartOfMessage byte
	MsgType        byte
	AuthRec        types.AuthRec
}

func NewGetAuthRec(msg []byte) (*GetAuthRec, error) {
	serialNumber := binary.LittleEndian.Uint32(msg[4:8])
	records := binary.LittleEndian.Uint32(msg[8:12])

	return &GetAuthRec{
		msg[0],
		msg[1],
		types.AuthRec{
			serialNumber,
			records,
		}}, nil
}
