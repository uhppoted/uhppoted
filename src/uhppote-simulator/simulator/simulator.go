package simulator

import (
	"errors"
	"fmt"
	"time"
	"uhppote/messages"
)

type Simulator struct {
	Debug bool
}

func (s *Simulator) Handle(bytes []byte) ([]byte, error) {
	if len(bytes) != 64 {
		return []byte{}, errors.New(fmt.Sprintf("Invalid message length %d", len(bytes)))
	}

	if bytes[0] != 0x17 {
		return []byte{}, errors.New(fmt.Sprintf("Invalid message type %02X", bytes[0]))
	}

	switch bytes[1] {
	case 0x94:
		return s.search(bytes)
	default:
		return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
	}

	return []byte{}, errors.New(fmt.Sprintf("Invalid command %02X", bytes[1]))
}

func (s *Simulator) search(bytes []byte) ([]byte, error) {
	time.Sleep(100 * time.Millisecond)

	msg := messages.Search{}

	return msg.Encode()
}
