package messages

import (
	"uhppote/types"
)

type GetStatus struct {
	StartOfMessage byte
	MsgType        byte
	Status         types.Status
}

func NewGetStatus(msg []byte) (*GetStatus, error) {
	status, err := types.DecodeStatus(msg)
	if err != nil {
		return nil, err
	}

	return &GetStatus{
		StartOfMessage: msg[0],
		MsgType:        msg[1],
		Status:         *status}, nil
}
