package auth

import (
	"github.com/uhppoted/uhppoted-api/kvs"
	"log"
	"strconv"
	"testing"
)

func TestValidateHOTPWithValidOTP(t *testing.T) {
	hotp := HOTP{
		increment: 8,
		secrets:   kvs.NewKeyValueStore("test:secrets", func(v string) (interface{}, error) { return v, nil }),
		counters: struct {
			*kvs.KeyValueStore
			filepath string
			log      *log.Logger
		}{
			kvs.NewKeyValueStore("test:counters", func(v string) (interface{}, error) { return strconv.ParseUint(v, 10, 64) }),
			"",
			nil,
		},
	}

	hotp.secrets.Put("qwerty", "DFIOJ3BJPHPCRJBT")

	if err := hotp.Validate("qwerty", "644039"); err != nil {
		t.Errorf("HOTP refused valid OTP")
	}

	if err := hotp.Validate("qwerty", "586787"); err != nil {
		t.Errorf("HOTP refused valid OTP")
	}
}

func TestValidateHOTPWithOutOfOrderOTP(t *testing.T) {
	hotp := HOTP{
		increment: 8,
		secrets:   kvs.NewKeyValueStore("test:secrets", func(v string) (interface{}, error) { return v, nil }),
		counters: struct {
			*kvs.KeyValueStore
			filepath string
			log      *log.Logger
		}{
			kvs.NewKeyValueStore("test:counters", func(v string) (interface{}, error) { return strconv.ParseUint(v, 10, 64) }),
			"",
			nil,
		},
	}

	hotp.secrets.Put("qwerty", "DFIOJ3BJPHPCRJBT")

	if err := hotp.Validate("qwerty", "586787"); err != nil {
		t.Errorf("HOTP refused valid OTP")
	}

	if err := hotp.Validate("qwerty", "644039"); err == nil {
		t.Errorf("HOTP accepted out of order OTP")
	}
}

func TestValidateHOTPWithOutOfRangeOTP(t *testing.T) {
	hotp := HOTP{
		increment: 2,
		secrets:   kvs.NewKeyValueStore("test:secrets", func(v string) (interface{}, error) { return v, nil }),
		counters: struct {
			*kvs.KeyValueStore
			filepath string
			log      *log.Logger
		}{
			kvs.NewKeyValueStore("test:counters", func(v string) (interface{}, error) { return strconv.ParseUint(v, 10, 64) }),
			"",
			nil,
		},
	}

	hotp.secrets.Put("qwerty", "DFIOJ3BJPHPCRJBT")

	if err := hotp.Validate("qwerty", "586787"); err == nil {
		t.Errorf("HOTP accepted out of range OTP")
	}
}

func TestValidateHOTPWithInvalidOTP(t *testing.T) {
	hotp := HOTP{
		increment: 8,
		secrets:   kvs.NewKeyValueStore("test:secrets", func(v string) (interface{}, error) { return v, nil }),
		counters: struct {
			*kvs.KeyValueStore
			filepath string
			log      *log.Logger
		}{
			kvs.NewKeyValueStore("test:counters", func(v string) (interface{}, error) { return strconv.ParseUint(v, 10, 64) }),
			"",
			nil,
		},
	}

	hotp.secrets.Put("qwerty", "DFIOJ3BJPHPCRJBT")

	if err := hotp.Validate("qwerty", "644038"); err == nil {
		t.Errorf("HOTP accepted invalid OTP")
	}
}
