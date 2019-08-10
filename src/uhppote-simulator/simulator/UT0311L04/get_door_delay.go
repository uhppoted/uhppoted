package UT0311L04

import (
	"net"
	"uhppote/messages"
)

func (s *UT0311L04) getDoorDelay(addr *net.UDPAddr, request *messages.GetDoorDelayRequest) {
	if request.SerialNumber == s.SerialNumber {

		if !(request.Door < 1 || request.Door > 4) {
			response := messages.GetDoorDelayResponse{
				SerialNumber: s.SerialNumber,
				Door:         request.Door,
				Unit:         0x03,
				Delay:        s.Doors[request.Door].Delay.Seconds(),
			}

			s.send(addr, &response)
		}
	}
}
