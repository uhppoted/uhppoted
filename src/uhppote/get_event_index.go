package uhppote

import (
	"uhppote/types"
)

type GetEventIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb4"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetEventIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb4"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetEventIndex(serialNumber uint32) (*types.EventIndex, error) {
	request := GetEventIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetEventIndexResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.EventIndex{
		SerialNumber: reply.SerialNumber,
		Index:        reply.Index,
	}, nil
}
