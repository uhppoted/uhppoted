package auth

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func load(filepath string, g func(key, value string) error) error {
	if filepath == "" {
		return nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer f.Close()

	re := regexp.MustCompile(`^\s*(.*?)\s+(\S.*)\s*`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])
			if err = g(key, value); err != nil {
				return err
			}
		}
	}

	return s.Err()
}

func store(filepath string, kv map[string]uint64) error {
	if filepath == "" {
		return nil
	}

	dir := path.Dir(filepath)
	filename := path.Base(filepath) + ".tmp"
	tmpfile := path.Join(dir, filename)

	f, err := os.Create(tmpfile)
	if err != nil {
		return err
	}

	defer f.Close()

	for key, value := range kv {
		if _, err := fmt.Fprintf(f, "%-20s %v\n", key, value); err != nil {
			return err
		}
	}

	f.Close()

	return os.Rename(tmpfile, filepath)
}

// NOTE: interim file watcher implementation pending fsnotify in Go 1.4
func watch(filepath string, reload func() error, logger *log.Logger) {
	go func() {
		finfo, err := os.Stat(filepath)
		if err != nil {
			log.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
			return
		}

		lastModified := finfo.ModTime()
		logged := false
		for {
			time.Sleep(2500 * time.Millisecond)
			finfo, err := os.Stat(filepath)
			if err != nil {
				if !logged {
					log.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
					logged = true
				}
			} else {
				logged = false
				if finfo.ModTime() != lastModified {
					log.Printf("INFO  Reloading information from %s\n", filepath)
					if err := reload(); err != nil {
						log.Printf("ERROR Failed to reload information from '%s': %v", filepath, err)
					} else {
						lastModified = finfo.ModTime()
					}
				}
			}
		}
	}()
}
