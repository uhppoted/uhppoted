package uhppote

import (
	"uhppote/messages"
)

func Marshal(m messages.Message) (*[]byte, error) {
	bytes := make([]byte, 64)

	bytes[0] = 0x17
	bytes[1] = m.Code()

	return &bytes, nil
}

func Unmarshal(bytes []byte) (string, error) {
	return "", nil
}
