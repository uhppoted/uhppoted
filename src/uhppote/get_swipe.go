package uhppote

import (
	"encoding/binary"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetSwipe(index uint32) (*types.Swipe, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0xb0
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], u.SerialNumber)
	binary.LittleEndian.PutUint32(cmd[8:12], index)

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewGetSwipe(reply)

	if err != nil {
		return nil, err
	}

	return result.Swipe, nil
}
