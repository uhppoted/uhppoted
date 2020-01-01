package auth

import (
	"fmt"
)

type RSA struct {
}

func NewRSA() (*RSA, error) {
	rsa := RSA{}

	return &rsa, nil
}

func (rsa *RSA) Validate(clientID string, request []byte, signature []byte) error {
	return fmt.Errorf("%s: invalid signature", clientID)
}
