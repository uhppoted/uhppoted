package UTC0311L04

import (
	"fmt"
	"net"
	"uhppote-simulator/entities"
	"uhppote/messages"
)

func (s *UTC0311L04) setDoorDelay(addr *net.UDPAddr, request *messages.SetDoorDelayRequest) {
	if request.SerialNumber == s.SerialNumber {
		if request.Unit != 0x03 {
			fmt.Printf("ERROR: Invalid time unit' - expected: %02x, received:%02x", 0x03, request.Unit)
			return
		}

		door := request.Door
		if door < 1 || door > 4 {
			fmt.Printf("ERROR: Invalid door' - expected: [1..4], received:%d", request.Door)
			return
		}

		s.Doors[door].Delay = entities.Delay(uint64(request.Delay) * 1000000000)

		response := messages.SetDoorDelayResponse{
			SerialNumber: s.SerialNumber,
			Door:         door,
			Unit:         0x03,
			Delay:        s.Doors[door].Delay.Seconds(),
		}

		s.send(addr, &response)

		if err := s.Save(); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}
