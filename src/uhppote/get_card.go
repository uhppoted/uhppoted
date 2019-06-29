package uhppote

import (
	"errors"
	"fmt"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

type GetCardByIndexRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x5c"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Index        uint32             `uhppote:"offset:8"`
}

type GetCardByIdRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x5a"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
}

type GetCardByIndexResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x5c"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         types.Date         `uhppote:"offset:12"`
	To           types.Date         `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

type GetCardByIdResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x5a"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
	From         *types.Date        `uhppote:"offset:12"`
	To           *types.Date        `uhppote:"offset:16"`
	Door1        bool               `uhppote:"offset:20"`
	Door2        bool               `uhppote:"offset:21"`
	Door3        bool               `uhppote:"offset:22"`
	Door4        bool               `uhppote:"offset:23"`
}

func (u *UHPPOTE) GetCardByIndex(serialNumber, index uint32) (*types.Card, error) {
	request := GetCardByIndexRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Index:        index,
	}

	reply, err := u.Send(serialNumber, request)
	if err != nil {
		return nil, err
	}

	response := GetCardByIdResponse{}
	err = codec.Unmarshal(reply, &response)
	if err != nil {
		return nil, err
	}

	if uint32(response.SerialNumber) != serialNumber {
		return nil, errors.New(fmt.Sprintf("Incorrect serial number in response - expect '%v', received '%v'", serialNumber, response.SerialNumber))
	}

	if response.CardNumber == 0 {
		return nil, nil
	}

	if response.From == nil {
		return nil, errors.New(fmt.Sprintf("Invalid 'from' date in response"))
	}

	if response.To == nil {
		return nil, errors.New(fmt.Sprintf("Invalid 'to' date in response"))
	}

	return &types.Card{
		CardNumber: response.CardNumber,
		From:       *response.From,
		To:         *response.To,
		Doors:      []bool{response.Door1, response.Door2, response.Door3, response.Door4},
	}, nil
}

func (u *UHPPOTE) GetCardById(serialNumber, cardNumber uint32) (*types.Card, error) {
	request := GetCardByIdRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		CardNumber:   cardNumber,
	}

	reply, err := u.Send(serialNumber, request)
	if err != nil {
		return nil, err
	}

	response := GetCardByIdResponse{}
	err = codec.Unmarshal(reply, &response)
	if err != nil {
		return nil, err
	}

	if uint32(response.SerialNumber) != serialNumber {
		return nil, errors.New(fmt.Sprintf("Incorrect serial number in response - expect '%v', received '%v'", serialNumber, response.SerialNumber))
	}

	if response.CardNumber == 0 {
		return nil, nil
	}

	if response.CardNumber != cardNumber {
		return nil, errors.New(fmt.Sprintf("Incorrect card number in response - expect '%v', received '%v'", cardNumber, response.CardNumber))
	}

	if response.From == nil {
		return nil, errors.New(fmt.Sprintf("Invalid 'from' date in response"))
	}

	if response.To == nil {
		return nil, errors.New(fmt.Sprintf("Invalid 'to' date in response"))
	}

	return &types.Card{
		CardNumber: response.CardNumber,
		From:       *response.From,
		To:         *response.To,
		Doors:      []bool{response.Door1, response.Door2, response.Door3, response.Door4},
	}, nil
}
