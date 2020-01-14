package auth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"uhppoted/kvs"
)

type Nonce struct {
	ignore   bool
	counters struct {
		*kvs.KeyValueStore
		filepath string
	}
	log *log.Logger
}

func NewNonce(verify bool, filepath string, logger *log.Logger) (*Nonce, error) {
	var err error

	var f = func(value string) (interface{}, error) {
		return strconv.ParseUint(value, 10, 64)
	}

	nonce := Nonce{
		ignore: !verify,
		counters: struct {
			*kvs.KeyValueStore
			filepath string
		}{
			kvs.NewKeyValueStore("nonce:counters", f),
			filepath,
		},
		log: logger,
	}

	if err = nonce.counters.LoadFromFile(filepath); err != nil {
		log.Printf("WARN: %v", err)
	}

	return &nonce, nil
}

func (n *Nonce) Validate(clientID *string, nonce *uint64) error {
	if !n.ignore || (clientID != nil && nonce != nil) {
		if clientID == nil {
			return errors.New("missing 'client-id'")
		}

		if nonce == nil {
			return errors.New("missing 'nonce'")
		}

		c, ok := n.counters.Get(*clientID)
		if !ok {
			c = uint64(0)
		}

		if *nonce <= c.(uint64) {
			return fmt.Errorf("reused: %s, %d", *clientID, *nonce)
		}

		n.counters.Store(*clientID, *nonce, n.counters.filepath, n.log)
	}

	return nil
}
