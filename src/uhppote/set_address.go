package uhppote

import (
	"errors"
	"fmt"
	"net"
	"uhppote/types"
)

type SetAddressRequest struct {
	MsgType      types.MsgType      `uhppote:"value:0x96"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Address      net.IP             `uhppote:"offset:8"`
	Mask         net.IP             `uhppote:"offset:12"`
	Gateway      net.IP             `uhppote:"offset:16"`
	MagicNumber  uint32             `uhppote:"offset:20"`
}

type SetAddressResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x96"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	Succeeded    bool               `uhppote:"offset:8"`
}

func (u *UHPPOTE) SetAddress(serialNumber uint32, address, mask, gateway net.IP) (*types.Result, error) {
	if address.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP address: %v", address))
	}

	if mask.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid subnet mask: %v", mask))
	}

	if gateway.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid gateway address: %v", gateway))
	}

	request := SetAddressRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Address:      address,
		Mask:         mask,
		Gateway:      gateway,
		MagicNumber:  0x55aaaa55,
	}

	reply := SetAddressResponse{}

	err := u.Execute(request, &reply)
	if err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: reply.SerialNumber,
		Succeeded:    reply.Succeeded,
	}, nil
}
