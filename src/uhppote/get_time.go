package uhppote

import (
	"encoding/binary"
	"fmt"
	"uhppote/messages"
	"uhppote/types"
)

func GetTime(serialNumber uint32, debug bool) (*types.DateTime, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x32
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

	reply, err := Execute(cmd, debug)

	if err != nil {
		return nil, err
	}

	result, err := messages.NewGetTime(reply)

	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf(" ... %v\n", *result)
	}

	return &result.DateTime, nil
}
