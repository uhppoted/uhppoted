package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

func (u *UHPPOTE) GetCardCount(serialNumber uint32) (*types.RecordCount, error) {
	request := struct {
		MsgType      byte   `uhppote:"offset:1"`
		SerialNumber uint32 `uhppote:"offset:4"`
	}{
		0x58,
		serialNumber,
	}

	reply := struct {
		MsgType      byte   `uhppote:"offset:1"`
		SerialNumber uint32 `uhppote:"offset:4"`
		Records      uint32 `uhppote:"offset:8"`
	}{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0x58 {
		return nil, errors.New(fmt.Sprintf("GetCardCount returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.RecordCount{
		SerialNumber: serialNumber,
		Records:      reply.Records,
	}, nil
}
