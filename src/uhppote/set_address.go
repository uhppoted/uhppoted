package uhppote

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func (u *UHPPOTE) SetAddress(serialNumber uint32, address, mask, gateway net.IP) error {
	if address.To4() == nil {
		return errors.New(fmt.Sprintf("Invalid IP address: %v", address))
	}

	if mask.To4() == nil {
		return errors.New(fmt.Sprintf("Invalid subnet mask: %v", mask))
	}

	if gateway.To4() == nil {
		return errors.New(fmt.Sprintf("Invalid gateway address: %v", gateway))
	}

	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x96
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

	copy(cmd[8:12], address.To4())
	copy(cmd[12:16], mask.To4())
	copy(cmd[16:20], gateway.To4())

	cmd[20] = 0x55
	cmd[21] = 0xaa
	cmd[22] = 0xaa
	cmd[23] = 0x55

	_, err := u.Execute(cmd)

	return err
}
