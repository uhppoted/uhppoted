package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

type RSA struct {
	clientKeys map[string]*rsa.PublicKey
	log        *log.Logger
}

func NewRSA(keys string, logger *log.Logger) (*RSA, error) {
	rsa := RSA{
		clientKeys: map[string]*rsa.PublicKey{},
		log:        logger,
	}

	if err := rsa.load(keys); err != nil {
		log.Printf("WARN: %v", err)
	}

	return &rsa, nil
}

func (r *RSA) Validate(clientID string, request []byte, signature []byte) error {
	pubkey, ok := r.clientKeys[clientID]
	if !ok || pubkey == nil {
		return fmt.Errorf("%s: no RSA public key", clientID)
	}

	hash := sha256.Sum256(request)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("Invalid RSA signature (%v)", err)
	}

	return nil
}

func (r *RSA) load(dir string) error {
	filemode, err := os.Stat(dir)
	if err != nil {
		return err
	}

	if !filemode.IsDir() {
		return fmt.Errorf("%s is not a directory", dir)
	}

	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range list {
		filename := f.Name()
		clientID := strings.TrimSuffix(filename, path.Ext(filename))

		bytes, err := ioutil.ReadFile(path.Join(dir, filename))
		if err != nil {
			r.log.Printf("WARN: %v", err)
		}

		block, _ := pem.Decode(bytes)
		if block == nil || block.Type != "PUBLIC KEY" {
			r.log.Printf("WARN: %s is not a valid RSA public key", filename)
			continue
		}

		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			r.log.Printf("WARN: %s is not a valid RSA public key (%v)", filename, err)
			continue
		}

		pubkey, ok := key.(*rsa.PublicKey)
		if !ok {
			r.log.Printf("WARN: %s is not a valid RSA public key", filename)
			continue
		}

		r.clientKeys[clientID] = pubkey
	}

	return nil
}
