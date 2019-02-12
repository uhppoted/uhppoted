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
		b := byte(0)

		switch ch {
		case '0':
			b = 0
		case '1':
			b = 1
		case '2':
			b = 2
		case '3':
			b = 3
		case '4':
			b = 4
		case '5':
			b = 5
		case '6':
			b = 6
		case '7':
			b = 7
		case '8':
			b = 8
		case '9':
			b = 9
		default:
			return nil, errors.New(fmt.Sprintf("Invalid numeric string '%s'", s))
		}

		bytes[ix/2] *= 16
		bytes[ix/2] += b
		ix += 1
	}

	return &bytes, nil
}
