package uhppote

import (
	"errors"
	"fmt"
)

func makeErr(msg string, err error) error {
	return errors.New(fmt.Sprintf(msg+" [%v]", err))
}
