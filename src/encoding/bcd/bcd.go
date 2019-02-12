package bcd

import (
	"errors"
	"fmt"
)

func Encode(s string) (*[]byte, error) {
	N := (len(s) + 1) / 2
	bytes := make([]byte, N)
	ix := len(s) % 2

	for _, ch := range s {
		b := byte(0x00)

		switch ch {
		case '0':
			b = 0x00
		case '1':
			b = 0x01
		case '2':
			b = 0x02
		case '3':
			b = 0x03
		case '4':
			b = 0x04
		case '5':
			b = 0x05
		case '6':
			b = 0x06
		case '7':
			b = 0x07
		case '8':
			b = 0x08
		case '9':
			b = 0x09
		default:
			return nil, errors.New(fmt.Sprintf("Invalid numeric string '%s'", s))
		}

		bytes[ix/2] *= 16
		bytes[ix/2] += b
		ix += 1
	}

	return &bytes, nil
}
