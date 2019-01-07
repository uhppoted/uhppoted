package messages

import (
	"encoding/binary"
	"fmt"
	"net"
	"uhppote/types"
)

type Search struct {
	SOM     byte
	MsgType byte
	Device  types.Device
}

func NewSearch(msg []byte) (*Search, error) {
	serialNumber := binary.LittleEndian.Uint32(msg[4:8])
	ipAddress := net.IPv4(msg[8], msg[9], msg[10], msg[11])
	subnetMask := net.IPv4(msg[12], msg[13], msg[14], msg[15])
	gateway := net.IPv4(msg[16], msg[17], msg[18], msg[19])
	macAddress := []byte{msg[20], msg[21], msg[22], msg[23], msg[24], msg[25]}
	version := fmt.Sprintf("%02X%02X", msg[26], msg[27])
	date := fmt.Sprintf("%04X-%02X-%02X", msg[28:30], msg[30:31], msg[31:32])

	device := types.Device{serialNumber,
		ipAddress,
		subnetMask,
		gateway,
		macAddress,
		version,
		date,
	}

	return &Search{msg[0], msg[1], device}, nil
}
