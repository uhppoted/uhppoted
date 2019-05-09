package uhppote

import (
	"uhppote/types"
)

type SetEventIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0xb2"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
	MagicWord    uint32             `uhppote:"offset:12"`
}

type SetEventIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0xb2"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Success      bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetEventIndex(serialNumber, index uint32) (*types.EventIndexResult, error) {
	request := SetEventIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
		MagicWord:    0x55aaaa55,
	}

	reply := SetEventIndexResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.EventIndexResult{
		SerialNumber: reply.SerialNumber,
		Index:        index,
		Succeeded:    reply.Success,
	}, nil
}
