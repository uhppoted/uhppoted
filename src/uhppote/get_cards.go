package uhppote

import (
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
)

func (u *UHPPOTE) GetCards(serialNumber uint32) (*types.RecordCount, error) {
	request := messages.GetCardsRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := messages.GetCardsResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.RecordCount{
		SerialNumber: reply.SerialNumber,
		Records:      reply.Records,
	}, nil
}
