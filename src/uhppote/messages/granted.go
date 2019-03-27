package messages

import (
	"encoding/binary"
	"uhppote/types"
)

type AddAuth struct {
	StartOfMessage byte
	MsgType        byte
	Authorised     types.Authorised
}

func NewAddAuth(msg []byte) (*AddAuth, error) {
	serialNumber := binary.LittleEndian.Uint32(msg[4:8])
	authorised := msg[8] == 0x01

	return &AddAuth{msg[0], msg[1], types.Authorised{
		serialNumber,
		authorised,
	}}, nil
}
