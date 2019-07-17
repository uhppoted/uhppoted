package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetEventIndex(serialNumber, index uint32) (*types.EventIndexResult, error) {
	request := messages.SetEventIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
		MagicWord:    0x55aaaa55,
	}

	reply := messages.SetEventIndexResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.EventIndexResult{
		SerialNumber: reply.SerialNumber,
		Index:        index,
		Changed:      reply.Changed,
	}, nil
}
