package uhppote

import (
	"uhppote/types"
)

type GetEventRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb0"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

type GetEventResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb0"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	Type         uint8              `uhppote:"offset:12"`
	Granted      bool               `uhppote:"offset:13"`
	Door         uint8              `uhppote:"offset:14"`
	DoorOpened   bool               `uhppote:"offset:15"`
	UserId       uint32             `uhppote:"offset:16"`
	Timestamp    types.DateTime     `uhppote:"offset:20"`
	RecordType   uint8              `uhppote:"offset:27"`
}

func (u *UHPPOTE) GetEvent(serialNumber, index uint32) (*types.Event, error) {
	request := GetEventRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
	}

	reply := GetEventResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Event{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
		Type:         reply.Type,
		Granted:      reply.Granted,
		Door:         reply.Door,
		DoorOpened:   reply.DoorOpened,
		UserId:       reply.UserId,
		Timestamp:    reply.Timestamp,
		RecordType:   reply.RecordType,
	}, nil
}
