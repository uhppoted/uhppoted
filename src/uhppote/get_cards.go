package uhppote

import (
	"uhppote/types"
)

type GetCardsRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x58"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
}

type GetCardsResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x58"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Records      uint32             `uhppote:"offset:8"`
}

func (u *UHPPOTE) GetCards(serialNumber uint32) (*types.RecordCount, error) {
	request := GetCardsRequest{
		SerialNumber: types.SerialNumber(serialNumber),
	}

	reply := GetCardsResponse{}

	err := u.Execute(serialNumber, request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.RecordCount{
		SerialNumber: reply.SerialNumber,
		Records:      reply.Records,
	}, nil
}
