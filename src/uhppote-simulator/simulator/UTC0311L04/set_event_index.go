package UTC0311L04

import (
	"errors"
	"fmt"
	"net"
	"uhppote/messages"
)

func (s *UTC0311L04) setEventIndex(addr *net.UDPAddr, request *messages.SetEventIndexRequest) {
	if s.SerialNumber == request.SerialNumber {

		if request.MagicWord != 0x55aaaa55 {
			s.onError(errors.New(fmt.Sprintf("Invalid 'magic number' - expected: %08x, received:%08x", 0x55aaaa55, request.MagicWord)))
		} else {

			updated := s.Events.SetIndex(request.Index)

			response := messages.SetEventIndexResponse{
				SerialNumber: s.SerialNumber,
				Changed:      updated,
			}

			s.send(addr, &response)

			if updated {
				if err := s.Save(); err != nil {
					s.onError(err)
				}
			}

		}
	}
}
