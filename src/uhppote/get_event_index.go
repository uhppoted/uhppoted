package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) GetEventIndex(serialNumber uint32) (*types.EventIndex, error) {
	request := messages.GetEventIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := messages.GetEventIndexResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.EventIndex{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
	}, nil
}
