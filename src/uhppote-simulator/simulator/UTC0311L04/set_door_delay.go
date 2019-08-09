package UTC0311L04

import (
	"net"
	"uhppote-simulator/entities"
	"uhppote/messages"
)

func (s *UTC0311L04) setDoorDelay(addr *net.UDPAddr, request *messages.SetDoorDelayRequest) {
	if request.SerialNumber == s.SerialNumber && request.Unit == 0x03 {
		door := request.Door
		if !(door < 1 || door > 4) {
			s.Doors[door].Delay = entities.Delay(uint64(request.Delay) * 1000000000)

			response := messages.SetDoorDelayResponse{
				SerialNumber: s.SerialNumber,
				Door:         door,
				Unit:         0x03,
				Delay:        s.Doors[door].Delay.Seconds(),
			}

			s.send(addr, &response)

			err := s.Save()
			if err != nil {
				s.onError(err)
			}
		}
	}
}
