package uhppote

import (
	"fmt"
	"uhppote/messages"
	"uhppote/types"
)

func Search(debug bool) ([]types.Device, error) {
	devices := []types.Device{}
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x94
	cmd[2] = 0x00
	cmd[3] = 0x00

	reply, err := Execute(cmd, debug)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewSearch(reply)

	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf(" ... %v\n", *result)
	}

	devices = append(devices, result.Device)

	return devices, nil
}
