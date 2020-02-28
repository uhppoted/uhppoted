package UT0311L04

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/messages"
	"github.com/uhppoted/uhppoted/src/uhppote-simulator/entities"
	"net"
)

func (s *UT0311L04) setDoorControlState(addr *net.UDPAddr, request *messages.SetDoorControlStateRequest) {
	if request.SerialNumber == s.SerialNumber {
		door := request.Door
		if door < 1 || door > 4 {
			fmt.Printf("ERROR: Invalid door' - expected: [1..4], received:%d", request.Door)
			return
		}

		s.Doors[door].ControlState = request.ControlState
		s.Doors[door].Delay = entities.Delay(uint64(request.Delay) * 1000000000)

		response := messages.SetDoorControlStateResponse{
			SerialNumber: s.SerialNumber,
			Door:         door,
			ControlState: s.Doors[door].ControlState,
			Delay:        s.Doors[door].Delay.Seconds(),
		}

		s.send(addr, &response)

		if err := s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
