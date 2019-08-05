package simulator

import (
	"errors"
	"fmt"
	"net"
	"uhppote/messages"
)

func (s *Simulator) setAddress(addr *net.UDPAddr, request *messages.SetAddressRequest) {
	if s.SerialNumber == request.SerialNumber {
		if request.MagicWord == 0x55aaaa55 {
			s.IpAddress = request.Address
			s.SubnetMask = request.Mask
			s.Gateway = request.Gateway

			err := s.Save()
			if err != nil {
				s.onError(err)
			}
		} else {
			s.onError(errors.New(fmt.Sprintf("Invalid 'magic number' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)))
		}
	}
}
