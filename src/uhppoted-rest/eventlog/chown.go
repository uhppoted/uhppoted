// +build !linux

package eventlog

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
