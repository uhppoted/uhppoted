package uhppote

import (
	"encoding/binary"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetStatus(serialNumber uint32) (*types.Status, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x20
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewGetStatus(reply)

	if err != nil {
		return nil, err
	}

	return &result.Status, nil
}
