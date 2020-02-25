package uhppoted

import (
	"strings"
)

func IsDevNull(path string) bool {
	return strings.ToLower(strings.TrimSpace(path)) == "nul"
}
