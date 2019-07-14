package messages

import (
	"encoding/hex"
	"fmt"
	"regexp"
)

type Request interface {
}

type Response interface {
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
