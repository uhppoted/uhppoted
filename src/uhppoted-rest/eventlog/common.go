package eventlog

import (
	"os"
	"path/filepath"
	"time"
)

const (
	backupTimeFormat = "2006-01-02T15-04-05.000"
	defaultMaxSize   = 100
)

var (
	currentTime = time.Now
	stat        = os.Stat
	megabyte    = 1024 * 1024
)

func deleteAll(dir string, files []logInfo) {
	for _, f := range files {
		_ = os.Remove(filepath.Join(dir, f.Name()))
	}
}

type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
