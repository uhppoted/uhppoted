package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type HOTP struct {
	Enabled   bool
	increment uint64
	secrets   struct {
		secrets  map[string]string
		filepath string
		guard    sync.Mutex
	}
	counters struct {
		counters map[string]uint64
		filepath string
		guard    sync.Mutex
	}
}

const DIGITS = 6

func NewHOTP(enabled bool, increment uint64, secrets string, counters string, logger *log.Logger) (*HOTP, error) {
	hotp := HOTP{
		Enabled:   enabled,
		increment: increment,
		secrets: struct {
			secrets  map[string]string
			filepath string
			guard    sync.Mutex
		}{
			secrets:  map[string]string{},
			filepath: secrets,
			guard:    sync.Mutex{},
		},
		counters: struct {
			counters map[string]uint64
			filepath string
			guard    sync.Mutex
		}{
			counters: map[string]uint64{},
			filepath: counters,
			guard:    sync.Mutex{},
		},
	}

	if enabled {
		keys, err := getSecrets(secrets)
		if err != nil {
			return nil, err
		}

		ctrs, err := getCounters(counters)
		if err != nil {
			return nil, err
		}

		hotp.secrets.secrets = keys
		hotp.counters.counters = ctrs

		f := func() error {
			return hotp.reload(logger)
		}

		watch(secrets, f, logger)
	}

	return &hotp, nil
}

func (hotp *HOTP) Validate(clientID, otp string) error {
	otp = strings.TrimSpace(otp)
	if len(otp) != DIGITS {
		return fmt.Errorf("%s: invalid OTP '%s' - expected %d digits", clientID, otp, DIGITS)
	}

	hotp.secrets.guard.Lock()
	defer hotp.secrets.guard.Unlock()

	secret, ok := hotp.secrets.secrets[clientID]
	if !ok {
		return fmt.Errorf("%s: no authorisation key", clientID)
	}

	hotp.counters.guard.Lock()
	defer hotp.counters.guard.Unlock()

	counter, ok := hotp.counters.counters[clientID]
	if !ok {
		counter = 1
	}

	for i := uint64(0); i < hotp.increment; i++ {
		generated, err := generateHOTP(secret, counter, DIGITS, sha1.New)
		if err != nil {
			return err
		}

		if subtle.ConstantTimeCompare([]byte(generated), []byte(otp)) == 1 {
			hotp.counters.counters[clientID] = counter + 1
			err := store(hotp.counters.filepath, hotp.counters.counters)
			if err != nil {
				fmt.Printf("WARN: Error storing updated HOTP counters (%v)\n", err)
			}
			return nil
		}

		counter++
	}

	return fmt.Errorf("%s: invalid OTP %s", clientID, otp)
}

// Ref. https://github.com/pquerna/otp
func generateHOTP(secret string, counter uint64, digits int, algorithm func() hash.Hash) (passcode string, err error) {
	secret = strings.TrimSpace(secret)
	if n := len(secret) % 8; n != 0 {
		secret = secret + strings.Repeat("=", 8-n)
	}

	bytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "", err
	}

	buffer := make([]byte, 8)
	mac := hmac.New(algorithm, bytes)
	binary.BigEndian.PutUint64(buffer, counter)

	mac.Write(buffer)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	mod := int32(value % int64(math.Pow10(digits)))

	return fmt.Sprintf("%06d", mod), nil
}

func getSecrets(path string) (map[string]string, error) {
	secrets := map[string]string{}
	err := load(path, func(key, value string) error {
		secrets[key] = value
		return nil
	})

	return secrets, err
}

func getCounters(path string) (map[string]uint64, error) {
	counters := map[string]uint64{}
	err := load(path, func(key, value string) error {
		if v, e := strconv.ParseUint(value, 10, 64); e != nil {
			return fmt.Errorf("Error parsing %s: %v", path, e)
		} else {
			counters[key] = v
			return nil
		}
	})

	return counters, err
}

func (h *HOTP) reload(log *log.Logger) error {
	secrets, err := getSecrets(h.secrets.filepath)
	if err != nil {
		return err
	}

	h.secrets.guard.Lock()
	defer h.secrets.guard.Unlock()

	if !reflect.DeepEqual(secrets, h.secrets.secrets) {
		for k, v := range secrets {
			h.secrets.secrets[k] = v
		}

		for k, _ := range h.secrets.secrets {
			if _, ok := secrets[k]; !ok {
				delete(h.secrets.secrets, k)
			}
		}

		log.Printf("WARN  Updated HOTP secrets from '%s'", h.secrets.filepath)
	}

	return nil
}
