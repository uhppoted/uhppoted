package uhppote

import (
	"encoding/binary"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) Authorise(serialNumber, cardNumber uint32, from, to types.Date, doors []int) (*types.Authorised, error) {
	cmd := make([]byte, 64)
	permissions := make([]byte, 4)

	for _, door := range doors {
		switch door {
		case 1:
			permissions[0] = 0x01
		case 2:
			permissions[1] = 0x01
		case 3:
			permissions[2] = 0x01
		case 4:
			permissions[3] = 0x01
		}
	}

	cmd[0] = 0x17
	cmd[1] = 0x50
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)
	binary.LittleEndian.PutUint32(cmd[8:12], cardNumber)

	from.Encode(cmd[12:16])
	to.Encode(cmd[16:20])
	copy(cmd[20:24], permissions)

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewAddAuth(reply)

	if err != nil {
		return nil, err
	}

	return &result.Authorised, nil
}
