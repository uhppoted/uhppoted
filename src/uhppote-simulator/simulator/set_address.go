package simulator

import (
	"errors"
	"fmt"
	"uhppote"
)

func (s *Simulator) SetAddress(request *uhppote.SetAddressRequest) (interface{}, error) {
	if s.SerialNumber != request.SerialNumber {
		return nil, nil
	}

	if request.MagicNumber != 0x55aaaa55 {
		return nil, errors.New(fmt.Sprintf("Invalid 'magic number' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicNumber))
	}

	s.IpAddress = request.Address
	s.SubnetMask = request.Mask
	s.Gateway = request.Gateway

	err := s.Save()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
