package uhppote

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
)

func makeErr(msg string, err error) error {
	return errors.New(fmt.Sprintf(msg+" [%v]", err))
}

func print(m []byte) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), "$1"))
}
