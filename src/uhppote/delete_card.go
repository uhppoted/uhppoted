package uhppote

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Revoked struct {
	SerialNumber uint32
	CardNumber   uint32
	Revoked      bool
}

func (r *Revoked) String() string {
	return fmt.Sprintf("REVOKED %v %v %v", r.SerialNumber, r.CardNumber, r.Revoked)
}

func (u *UHPPOTE) Revoke(serialNumber, cardNumber uint32) (*Revoked, error) {
	cmd := encodeRevokeRequest(serialNumber, cardNumber)
	reply, err := u.Execute(cmd)

	if err != nil {
		return nil, err
	}

	if len(reply) != 64 {
		return nil, errors.New(fmt.Sprintf("Invalid reply length: %v", len(reply)))
	}

	if reply[0] != 0x17 {
		return nil, errors.New(fmt.Sprintf("Invalid reply start of message: %02X", reply[0]))
	}

	if reply[1] != 0x52 {
		return nil, errors.New(fmt.Sprintf("Invalid reply command code: %02X", reply[1]))
	}

	if binary.LittleEndian.Uint32(reply[4:8]) != serialNumber {
		return nil, errors.New(fmt.Sprintf("Invalid reply serial number: %v", binary.LittleEndian.Uint32(reply[4:8])))
	}

	if reply[8] != 0x00 && reply[8] != 0x01 {
		return nil, errors.New(fmt.Sprintf("Invalid reply result code: %02X", reply[8]))
	}

	return &Revoked{
		SerialNumber: binary.LittleEndian.Uint32(reply[4:8]),
		CardNumber:   cardNumber,
		Revoked:      reply[8] == 0x01,
	}, nil
}

func encodeRevokeRequest(serialNumber, cardNumber uint32) []byte {
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x52
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)
	binary.LittleEndian.PutUint32(cmd[8:12], cardNumber)

	return cmd
}