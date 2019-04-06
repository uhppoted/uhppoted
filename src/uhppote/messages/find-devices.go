package messages

import (
	"encoding/binary"
	"fmt"
	"net"
	"uhppote/types"
)

type FindDevicesResponse struct {
	SOM     byte
	MsgType byte
	Device  types.Device
}

func NewFindDevicesRequest() (*Message, error) {
	return &Message{0x94, []Field{}}, nil
}

func NewFindDevicesResponse(msg []byte) (*FindDevicesResponse, error) {
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

	return &FindDevicesResponse{msg[0], msg[1], device}, nil
}

func (s *FindDevicesResponse) Encode() ([]byte, error) {
	reply := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x01, 0x7d, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	return reply, nil
}
