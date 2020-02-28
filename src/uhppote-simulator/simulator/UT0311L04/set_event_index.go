package UT0311L04

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/messages"
	"net"
)

func (s *UT0311L04) setEventIndex(addr *net.UDPAddr, request *messages.SetEventIndexRequest) {
	if s.SerialNumber == request.SerialNumber {
		if request.MagicWord != 0x55aaaa55 {
			fmt.Printf("ERROR: Invalid 'magic word' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)
			return
		}

		updated := s.Events.SetIndex(request.Index)

		response := messages.SetEventIndexResponse{
			SerialNumber: s.SerialNumber,
			Changed:      updated,
		}

		s.send(addr, &response)

		if updated {
			if err := s.Save(); err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		}
	}
}
