package uhppote

import (
	"encoding/binary"
	"fmt"
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func SetTime(serialNumber uint32, datetime time.Time, debug bool) (*types.DateTime, error) {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x30
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

	cmd[8] = encode(datetime.Year() / 100)
	cmd[9] = encode(datetime.Year() % 100)
	cmd[10] = encode(int(datetime.Month()))
	cmd[11] = encode(datetime.Day())
	cmd[12] = encode(datetime.Hour())
	cmd[13] = encode(datetime.Minute())
	cmd[14] = encode(datetime.Second())

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

func encode(b int) byte {
	msb := b / 10
	lsb := b % 10

	return byte(msb*16 + lsb)
}
