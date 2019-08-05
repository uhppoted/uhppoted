package uhppote

import (
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) DeleteCards(serialNumber uint32) (*types.Result, error) {
	request := messages.DeleteCardsRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		MagicWord:    0x55aaaa55,
	}

	reply := messages.DeleteCardsResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
