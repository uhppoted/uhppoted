package uhppote

import (
	"errors"
	"fmt"
	"uhppote/types"
)

type GetEventRequest struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

type GetEventResponse struct {
	MsgType      byte               `uhppote:"offset:1"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	Type         uint8              `uhppote:"offset:12"`
	Granted      bool               `uhppote:"offset:13"`
	Door         uint8              `uhppote:"offset:14"`
	DoorState    uint8              `uhppote:"offset:15"`
	CardNumber   uint32             `uhppote:"offset:16"`
	Timestamp    types.DateTime     `uhppote:"offset:20"`
	RecordType   uint8              `uhppote:"offset:27"`
}

func (u *UHPPOTE) GetEvent(serialNumber, index uint32) (*types.Event, error) {
	request := GetEventRequest{
		MsgType:      0xb0,
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
	}

	reply := GetEventResponse{}

	err := u.Exec(request, &reply)
	if err != nil {
		return nil, err
	}

	if reply.MsgType != 0xb0 {
		return nil, errors.New(fmt.Sprintf("GetEvent returned incorrect message type: %02X\n", reply.MsgType))
	}

	return &types.Event{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
		Type:         reply.Type,
		Granted:      reply.Granted,
		Door:         reply.Door,
		DoorState:    reply.DoorState,
		CardNumber:   reply.CardNumber,
		Timestamp:    reply.Timestamp,
		RecordType:   reply.RecordType,
	}, nil
}
