package uhppote

import (
	"encoding/binary"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetTime(datetime types.DateTime) (*types.DateTime, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x30
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], u.SerialNumber)
	datetime.Encode(cmd[8:15])

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewGetTime(reply)

	if err != nil {
		return nil, err
	}

	return &result.DateTime, nil
}
