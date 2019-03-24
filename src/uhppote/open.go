package uhppote

import (
	"encoding/binary"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) OpenDoor(serialNumber uint32, door byte) (*types.Opened, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x40
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)
	cmd[8] = door

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewOpenDoorReply(reply)

	if err != nil {
		return nil, err
	}

	opened := types.Opened{
		SerialNumber: serialNumber,
		Door:         uint32(door),
		Opened:       result.Opened,
	}

	return &opened, nil
}
