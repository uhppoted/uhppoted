package simulator

import (
	"net"
	"uhppote-simulator/entities"
	"uhppote/messages"
)

type Simulator interface {
	DeviceID() uint32
	DeviceType() string
	FilePath() string
	SetTxQ(chan entities.Message)

	Handle(*net.UDPAddr, messages.Request)
	Save() error
	Delete() error

	Swipe(deviceId uint32, cardNumber uint32, door uint8) (bool, uint32)
}
