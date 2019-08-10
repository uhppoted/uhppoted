package UTC0311L04

import (
	"fmt"
	"net"
	"uhppote/messages"
)

func (s *UTC0311L04) setAddress(addr *net.UDPAddr, request *messages.SetAddressRequest) {
	if s.SerialNumber == request.SerialNumber {
		if request.MagicWord != 0x55aaaa55 {
			fmt.Printf("ERROR: Invalid 'magic word' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)
			return
		}

		s.IpAddress = request.Address
		s.SubnetMask = request.Mask
		s.Gateway = request.Gateway

		if err := s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
