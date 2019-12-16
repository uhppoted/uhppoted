package kvs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type KeyValueStore struct {
	name  string
	store map[string]interface{}
	guard sync.Mutex
	re    *regexp.Regexp
	f     func(string) (interface{}, error)
}

func NewKeyValueStore(name string, f func(string) (interface{}, error)) *KeyValueStore {
	return &KeyValueStore{
		name:  name,
		store: map[string]interface{}{},
		guard: sync.Mutex{},
		re:    regexp.MustCompile(`^\s*(.*?)(?:\s{2,})(\S.*)\s*`),
		f:     f,
	}
}

func (kv *KeyValueStore) Get(key string) (interface{}, bool) {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	value, ok := kv.store[key]

	return value, ok
}

func (kv *KeyValueStore) Put(key string, value interface{}) {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	kv.store[key] = value
}

func (kv *KeyValueStore) LoadFromFile(filepath string) error {
	store, err := load(filepath, kv.f)
	if err != nil {
		return err
	}

	return kv.merge(*store)
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

// NOTE: interim file watcher implementation pending fsnotify in Go 1.4
func (kv *KeyValueStore) Watch(filepath string, logger *log.Logger) {
	go func() {
		finfo, err := os.Stat(filepath)
		if err != nil {
			logger.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
			return
		}

		lastModified := finfo.ModTime()
		logged := false
		for {
			time.Sleep(2500 * time.Millisecond)
			finfo, err := os.Stat(filepath)
			if err != nil {
				if !logged {
					logger.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
					logged = true
				}

				continue
			}

			logged = false
			if finfo.ModTime() != lastModified {
				log.Printf("INFO  Reloading information from %s\n", filepath)

				store, err := load(filepath, kv.f)
				if err != nil {
					log.Printf("ERROR Failed to reload information from %s: %v", filepath, err)
					continue
				}

				err = kv.merge(*store)
				if err != nil {
					logger.Printf("ERROR Failed to reload information from %s: %v", filepath, err)
					continue
				}

				log.Printf("WARN  Updated %s from %s", kv.name, filepath)
				lastModified = finfo.ModTime()
			}
		}
	}()
}

func load(filepath string, g func(string) (interface{}, error)) (*map[string]interface{}, error) {
	if filepath == "" {
		return nil, nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	store := map[string]interface{}{}
	re := regexp.MustCompile(`^\s*(.*?)(?:\s{2,})(\S.*)\s*`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			if v, err := g(value); err != nil {
				return &store, err
			} else {
				store[key] = v
			}
		}
	}

	return &store, s.Err()
}

func (kv *KeyValueStore) merge(store map[string]interface{}) error {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	if !reflect.DeepEqual(store, kv.store) {
		for k, v := range store {
			kv.store[k] = v
		}

		for k, _ := range kv.store {
			if _, ok := store[k]; !ok {
				delete(kv.store, k)
			}
		}
	}

	return nil
}
