package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetEvent(serialNumber, index uint32) (*types.Event, error) {
	request := messages.GetEventRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
	}

	reply := messages.GetEventResponse{}

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
		UserID:       reply.UserID,
		Timestamp:    reply.Timestamp,
		Result:       reply.Result,
	}, nil
}
