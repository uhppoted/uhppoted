package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) Search() ([]types.Device, error) {
	devices := []types.Device{}
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x94
	cmd[2] = 0x00
	cmd[3] = 0x00

	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewSearch(reply)

	if err != nil {
		return nil, err
	}

	devices = append(devices, result.Device)

	return devices, nil
}
