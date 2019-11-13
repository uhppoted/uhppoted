package eventlog

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var _ io.WriteCloser = (*Ticker)(nil)

type Ticker struct {
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"maxsize" yaml:"maxsize"`
	MaxAge     int    `json:"maxage" yaml:"maxage"`
	MaxBackups int    `json:"maxbackups" yaml:"maxbackups"`
	LocalTime  bool   `json:"localtime" yaml:"localtime"`

	size int64
	file *os.File
	mu   sync.Mutex
}

func (l *Ticker) Write(p []byte) (n int, err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > l.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, l.max(),
		)
	}

	if l.file == nil {
		if err = l.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if l.size+writeLen > l.max() {
		if err := l.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = l.file.Write(p)
	l.size += int64(n)

	return n, err
}

func (l *Ticker) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.close()
}

func (l *Ticker) close() error {
	if l.file == nil {
		return nil
	}
	err := l.file.Close()
	l.file = nil
	return err
}

func (l *Ticker) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotate()
}

func (l *Ticker) rotate() error {
	if err := l.close(); err != nil {
		return err
	}

	if err := l.openNew(); err != nil {
		return err
	}

	return l.cleanup()
}

func (l *Ticker) openNew() error {
	err := os.MkdirAll(l.dir(), 0744)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := l.filename()
	mode := os.FileMode(0644)
	//	err = l.archive()
	//	if err != nil {
	//		return fmt.Errorf("Error archiving log file: %s", err)
	//	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}

	l.file = f
	l.size = 0

	return nil
}

func (l *Ticker) archive() error {
	name := l.filename()
	info, err := os.Stat(name)

	if err == nil {
		newname := backupName(name, l.LocalTime)
		if err := os.Rename(name, newname); err != nil {
			return err
		}

		if err := chown(name, info); err != nil {
			return err
		}

		if err := l.compress(newname); err != nil {
			log.Printf("Error compressing archive file: %v\n", err)
		} else if err := os.Remove(newname); err != nil {
			log.Printf("Error deleting log file '%v'  %v\n", newname, err)
		}
	}

	return nil
}

func (l *Ticker) compress(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	gzfile := filepath + ".gz"
	gz, err := os.Create(gzfile)
	if err != nil {
		return err
	}

	defer gz.Close()

	w, err := gzip.NewWriterLevel(gz, 9)
	if err != nil {
		return err
	}

	_, err = w.Write(s)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := currentTime()
	if !local {
		t = t.UTC()
	}

	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

func (l *Ticker) openExistingOrNew(writeLen int) error {
	filename := l.filename()
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return l.openNew()
	}

	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= l.max() {
		return l.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return l.openNew()
	}

	l.file = file
	l.size = info.Size()

	return nil
}

func (l *Ticker) filename() string {
	if l.Filename != "" {
		return l.Filename
	}
	name := filepath.Base(os.Args[0]) + "-lumberjack.log"
	return filepath.Join(os.TempDir(), name)
}

func (l *Ticker) cleanup() error {
	if l.MaxBackups == 0 && l.MaxAge == 0 {
		return nil
	}

	files, err := l.oldLogFiles()
	if err != nil {
		return err
	}

	var deletes []logInfo

	if l.MaxBackups > 0 && l.MaxBackups < len(files) {
		deletes = files[l.MaxBackups:]
		files = files[:l.MaxBackups]
	}
	if l.MaxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(l.MaxAge))

		cutoff := currentTime().Add(-1 * diff)

		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				deletes = append(deletes, f)
			}
		}
	}

	if len(deletes) == 0 {
		return nil
	}

	go deleteAll(l.dir(), deletes)

	return nil
}

func (l *Ticker) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(l.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	logFiles := []logInfo{}

	prefix, ext := l.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := l.timeFromName(f.Name(), prefix, ext)
		if name == "" {
			continue
		}
		t, err := time.Parse(backupTimeFormat, name)
		if err == nil {
			logFiles = append(logFiles, logInfo{t, f})
		}
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

func (l *Ticker) timeFromName(filename, prefix, ext string) string {
	if !strings.HasPrefix(filename, prefix) {
		return ""
	}
	filename = filename[len(prefix):]

	if !strings.HasSuffix(filename, ext) {
		return ""
	}
	filename = filename[:len(filename)-len(ext)]
	return filename
}

func (l *Ticker) max() int64 {
	if l.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(l.MaxSize) * int64(megabyte)
}

func (l *Ticker) dir() string {
	return filepath.Dir(l.filename())
}

func (l *Ticker) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(l.filename())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}
