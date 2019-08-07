package simulator

import (
	"net"
	"time"
	"uhppote-simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

func (s *Simulator) openDoor(addr *net.UDPAddr, request *messages.OpenDoorRequest) {
	if s.SerialNumber == request.SerialNumber {
		granted := false
		opened := false
		door := request.Door

		if !(door < 1 || door > 4) {
			granted = true
			opened = s.Doors[door].Open()

			response := messages.OpenDoorResponse{
				SerialNumber: s.SerialNumber,
				Succeeded:    granted && opened,
			}

			s.send(addr, &response)

			datetime := time.Now().UTC().Add(time.Duration(s.TimeOffset))
			event := entities.Event{
				Type:       0x02,
				Granted:    granted,
				Door:       door,
				DoorOpened: opened,
				UserId:     3922570474,
				Timestamp:  types.DateTime(datetime),
				RecordType: 0x2c,
			}

			s.add(&event)
		}
	}
}
