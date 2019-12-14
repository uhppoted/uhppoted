package kvs

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type KeyValueStore struct {
	store map[string]interface{}
	re    *regexp.Regexp
	f     func(string) (interface{}, error)
}

func NewKeyValueStore(f func(string) (interface{}, error)) *KeyValueStore {
	return &KeyValueStore{
		store: map[string]interface{}{},
		re:    regexp.MustCompile(`^\s*(.*?)(?::\s*|\s*=\s*|\s{2,})(\S.*)\s*`),
		f:     f,
	}
}

func (kv *KeyValueStore) Get(key string) (interface{}, bool) {
	value, ok := kv.store[key]

	return value, ok
}

func (kv *KeyValueStore) Put(key string, value interface{}) {
	kv.store[key] = value
}

func (kv *KeyValueStore) Load(r io.Reader) error {
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := kv.re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			if v, err := kv.f(value); err != nil {
				return err
			} else {
				kv.store[key] = v
			}
		}
	}

	return s.Err()
}

func (kv *KeyValueStore) Save(w io.Writer) error {
	for key, value := range kv.store {
		if _, err := fmt.Fprintf(w, "%-20s  %v\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

// // NOTE: interim file watcher implementation pending fsnotify in Go 1.4
// func watch(filepath string, reload func() error, logger *log.Logger) {
// 	go func() {
// 		finfo, err := os.Stat(filepath)
// 		if err != nil {
// 			log.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
// 			return
// 		}
//
// 		lastModified := finfo.ModTime()
// 		logged := false
// 		for {
// 			time.Sleep(2500 * time.Millisecond)
// 			finfo, err := os.Stat(filepath)
// 			if err != nil {
// 				if !logged {
// 					log.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
// 					logged = true
// 				}
// 			} else {
// 				logged = false
// 				if finfo.ModTime() != lastModified {
// 					log.Printf("INFO  Reloading information from %s\n", filepath)
// 					if err := reload(); err != nil {
// 						log.Printf("ERROR Failed to reload information from '%s': %v", filepath, err)
// 					} else {
// 						lastModified = finfo.ModTime()
// 					}
// 				}
// 			}
// 		}
// 	}()
// }
