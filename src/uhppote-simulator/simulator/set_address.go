package simulator

import (
	"errors"
	"fmt"
	"uhppote/messages"
)

func (s *Simulator) setAddress(request *messages.SetAddressRequest) interface{} {
	if s.SerialNumber != request.SerialNumber {
		return nil
	}

	if request.MagicWord != 0x55aaaa55 {
		s.onError(errors.New(fmt.Sprintf("Invalid 'magic number' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)))
		return nil
	}

	s.IpAddress = request.Address
	s.SubnetMask = request.Mask
	s.Gateway = request.Gateway

	err := s.Save()
	if err != nil {
		s.onError(err)
		return nil
	}

	return nil
}
