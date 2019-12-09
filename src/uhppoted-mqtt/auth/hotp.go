package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"strings"
)

type HOTP struct {
	Enabled  bool
	Range    uint64
	Secrets  map[string]string
	Counters map[string]uint64
}

const DIGITS = 6

func NewHOTP(enabled bool, limit uint64, secrets string, counters string) (*HOTP, error) {
	hotp := HOTP{
		Enabled:  enabled,
		Range:    limit,
		Secrets:  map[string]string{},
		Counters: map[string]uint64{},
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

		hotp.Secrets = keys
		hotp.Counters = ctrs
	}

	return &hotp, nil
}

func getSecrets(path string) (map[string]string, error) {
	return map[string]string{
		"QWERTY54": "DFIOJ3BJPHPCRJBT",
	}, nil
}

func getCounters(path string) (map[string]uint64, error) {
	return map[string]uint64{
		"QWERTY54": 1,
	}, nil
}

func ValidateHOTP(passcode string, counter uint64, secret string) bool {
	passcode = strings.TrimSpace(passcode)
	if len(passcode) != DIGITS {
		return false
	}

	otpstr, err := generateHOTP(secret, counter, DIGITS, sha1.New)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(otpstr), []byte(passcode)) == 1
}

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
