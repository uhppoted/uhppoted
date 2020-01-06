package auth

import (
	"crypto/hmac"
	"crypto/sha256"
)

type HMAC struct {
	Required bool
	secret   []byte
}

func NewHMAC(required bool, key string) (*HMAC, error) {
	return &HMAC{
		Required: required,
		secret:   []byte(key),
	}, nil
}

func (h *HMAC) Verify(message []byte, mac []byte) bool {
	hash := hmac.New(sha256.New, h.secret)
	hash.Write(message)

	return hmac.Equal(mac, hash.Sum(nil))
}
