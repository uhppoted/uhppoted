package auth

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppoted/kvs"
	"log"
	"strconv"
)

type Nonce struct {
	ignore bool
	mqttd  struct {
		*kvs.KeyValueStore
		filepath string
	}
	counters struct {
		*kvs.KeyValueStore
		filepath string
	}
	log *log.Logger
}

func NewNonce(verify bool, server, clients string, logger *log.Logger) (*Nonce, error) {
	var err error

	var f = func(value string) (interface{}, error) {
		return strconv.ParseUint(value, 10, 64)
	}

	nonce := Nonce{
		ignore: !verify,
		mqttd: struct {
			*kvs.KeyValueStore
			filepath string
		}{
			kvs.NewKeyValueStore("nonce:mqttd", f),
			server,
		},
		counters: struct {
			*kvs.KeyValueStore
			filepath string
		}{
			kvs.NewKeyValueStore("nonce:clients", f),
			clients,
		},
		log: logger,
	}

	if err = nonce.mqttd.LoadFromFile(server); err != nil {
		log.Printf("WARN  %v", err)
	}

	if err = nonce.counters.LoadFromFile(clients); err != nil {
		log.Printf("WARN  %v", err)
	}

	return &nonce, nil
}

func (n *Nonce) Validate(clientID *string, nonce *uint64) error {
	if !n.ignore || (clientID != nil && nonce != nil) {
		if clientID == nil {
			return errors.New("missing client-id")
		}

		if nonce == nil {
			return errors.New("missing nonce missing")
		}

		c, ok := n.counters.Get(*clientID)
		if !ok {
			c = uint64(0)
		}

		if *nonce <= c.(uint64) {
			return fmt.Errorf("nonce reused: %s, %d", *clientID, *nonce)
		}

		n.counters.Store(*clientID, *nonce, n.counters.filepath, n.log)
	}

	return nil
}

func (n *Nonce) Next() uint64 {
	c, ok := n.mqttd.Get("mqttd")
	if !ok {
		c = uint64(0)
	}

	nonce := c.(uint64) + 1

	n.mqttd.Store("mqttd", nonce, n.mqttd.filepath, n.log)

	return nonce
}
