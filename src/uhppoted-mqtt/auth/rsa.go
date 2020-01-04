package auth

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"uhppoted/kvs"
)

type RSA struct {
	key        *rsa.PrivateKey
	clientKeys map[string]*rsa.PublicKey
	counters   struct {
		*kvs.KeyValueStore
		filepath string
	}
	log *log.Logger
}

func NewRSA(privateKey, clientKeys, counters string, logger *log.Logger) (*RSA, error) {
	rsa := RSA{
		key:        nil,
		clientKeys: map[string]*rsa.PublicKey{},
		counters: struct {
			*kvs.KeyValueStore
			filepath string
		}{
			kvs.NewKeyValueStore("rsa:counters", func(value string) (interface{}, error) { return strconv.ParseUint(value, 10, 64) }),
			counters,
		},
		log: logger,
	}

	if err := rsa.loadPrivateKey(privateKey); err != nil {
		log.Printf("WARN: %v", err)
	}

	if err := rsa.loadClientKeys(clientKeys); err != nil {
		log.Printf("WARN: %v", err)
	}

	if err := rsa.counters.LoadFromFile(counters); err != nil {
		log.Printf("WARN: %v", err)
	}

	// TODO 'watch' client key directory

	return &rsa, nil
}

func (r *RSA) Validate(clientID string, request []byte, signature []byte, counter uint64) error {
	pubkey, ok := r.clientKeys[clientID]
	if !ok || pubkey == nil {
		return fmt.Errorf("%s: no RSA public key", clientID)
	}

	hash := sha256.Sum256(request)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("%s: invalid RSA signature (%v)", clientID, err)
	}

	c, ok := r.counters.Get(clientID)
	if !ok {
		c = uint64(0)
	}

	if counter <= c.(uint64) {
		// r.log.Printf("TODO:  %s: RSA counter reused (%d)", clientID, counter)
		return fmt.Errorf("%s: RSA counter reused (%d)", clientID, counter)
	}

	r.counters.Store(clientID, counter, r.counters.filepath, r.log)

	return nil
}

func (r *RSA) Decrypt(ciphertext []byte, iv []byte, key []byte) ([]byte, error) {
	secretKey, err := rsa.DecryptPKCS1v15(nil, r.key, key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext missing IV (%d bytes)", len(ciphertext))
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext not a multiple of AES block size (%d bytes)", len(ciphertext))
	}

	plaintext := make([]byte, len(ciphertext))

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	N := len(plaintext)
	padding := int(plaintext[N-1])
	N -= padding

	if N < 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	return plaintext[:N], nil
}

func (r *RSA) loadPrivateKey(filepath string) error {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	r.key = pk

	return nil
}
func (r *RSA) loadClientKeys(dir string) error {
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
