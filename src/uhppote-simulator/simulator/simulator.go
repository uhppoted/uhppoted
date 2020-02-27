package simulator

import (
	"github.com/uhppoted/uhppoted/src/uhppote-simulator/entities"
	"github.com/uhppoted/uhppoted/src/uhppote/messages"
	"net"
)

type Simulator interface {
	DeviceID() uint32
	DeviceType() string
	FilePath() string
	SetTxQ(chan entities.Message)

	Handle(*net.UDPAddr, messages.Request)
	Save() error
	Delete() error

	Swipe(deviceID uint32, cardNumber uint32, door uint8) (bool, uint32)
}
