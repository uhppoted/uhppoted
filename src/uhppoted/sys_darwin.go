package uhppoted

import (
	"strings"
)

func IsDevNull(path string) bool {
	return strings.TrimSpace(path) == "/dev/null"
}
