package auth

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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

type keyset struct {
	key        *rsa.PrivateKey
	clientKeys map[string]*rsa.PublicKey
}

type RSA struct {
	signingKeys    keyset
	encryptionKeys keyset
	log            *log.Logger
}

func NewRSA(keydir string, logger *log.Logger) (*RSA, error) {
	var err error

	r := RSA{
		signingKeys: keyset{
			key:        nil,
			clientKeys: map[string]*rsa.PublicKey{},
		},
		encryptionKeys: keyset{
			key:        nil,
			clientKeys: map[string]*rsa.PublicKey{},
		},
		log: logger,
	}

	r.signingKeys.key, err = loadPrivateKey(path.Join(keydir, "signing", "private.key"))
	if err != nil {
		log.Printf("WARN: %v", err)
	}

	r.signingKeys.clientKeys, err = loadPublicKeys(path.Join(keydir, "signing"), logger)
	if err != nil {
		log.Printf("WARN: %v", err)
	}

	r.encryptionKeys.key, err = loadPrivateKey(path.Join(keydir, "encryption", "private.key"))
	if err != nil {
		log.Printf("WARN: %v", err)
	}

	r.encryptionKeys.clientKeys, err = loadPublicKeys(path.Join(keydir, "encryption"), logger)
	if err != nil {
		log.Printf("WARN: %v", err)
	}

	// TODO 'watch' client key directory

	return &r, nil
}

func (r *RSA) Validate(clientID string, request []byte, signature []byte) error {
	pubkey, ok := r.signingKeys.clientKeys[clientID]
	if !ok || pubkey == nil {
		return fmt.Errorf("%s: no RSA public key", clientID)
	}

	hash := sha256.Sum256(request)
	err := rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("%s: invalid RSA signature (%v)", clientID, err)
	}

	return nil
}

func (r *RSA) Sign(message []byte) ([]byte, error) {
	key := r.signingKeys.key
	if key != nil {
		rng := rand.Reader
		hashed := sha256.Sum256(message)

		return rsa.SignPKCS1v15(rng, key, crypto.SHA256, hashed[:])
	}

	return []byte{}, nil
}

func (r *RSA) Encrypt(plaintext []byte, clientID string) ([]byte, []byte, []byte, error) {
	secretKey := make([]byte, 32)
	if _, err := rand.Read(secretKey); err != nil {
		return nil, nil, nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, nil, err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, nil, nil, err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	ciphertext := make([]byte, len(plaintext)+padding)
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext, append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...))

	rng := rand.Reader
	label := []byte{}
	key, err := rsa.EncryptOAEP(sha256.New(), rng, r.encryptionKeys.clientKeys[clientID], secretKey, label)
	if err != nil {
		return nil, nil, nil, err
	}

	return ciphertext, iv, key, nil
}

func (r *RSA) Decrypt(ciphertext []byte, iv []byte, key []byte) ([]byte, error) {
	secretKey, err := rsa.DecryptPKCS1v15(nil, r.encryptionKeys.key, key)
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

	// Shouldn't really need this but using openssl AES on the command with the -salt option prepends the
	// actual ciphertext with 'Salted__<salt>'
	// Ref. http://justsolve.archiveteam.org/wiki/OpenSSL_salted_format
	offset := 0
	if strings.HasPrefix(string(ciphertext), "Salted__") {
		offset = 16
	}

	plaintext := make([]byte, len(ciphertext[offset:]))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext[offset:])

	N := len(plaintext)
	padding := int(plaintext[N-1])
	N -= padding

	if N < 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	return plaintext[:N], nil
}

func loadPrivateKey(filepath string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s is not a valid RSA private key", filepath)
	}

	return pk, nil
}

func loadPublicKeys(dir string, log *log.Logger) (map[string]*rsa.PublicKey, error) {
	keys := map[string]*rsa.PublicKey{}
	filemode, err := os.Stat(dir)
	if err != nil {
		return keys, err
	}

	if !filemode.IsDir() {
		return keys, fmt.Errorf("%s is not a directory", dir)
	}

	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return keys, err
	}

	for _, f := range list {
		filename := f.Name()
		ext := path.Ext(filename)
		if ext == ".pub" {
			clientID := strings.TrimSuffix(filename, ext)

			bytes, err := ioutil.ReadFile(path.Join(dir, filename))
			if err != nil {
				log.Printf("WARN: %v", err)
			}

			block, _ := pem.Decode(bytes)
			if block == nil || block.Type != "PUBLIC KEY" {
				log.Printf("WARN: %s is not a valid RSA public key", filename)
				continue
			}

			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				log.Printf("WARN: %s is not a valid RSA public key (%v)", filename, err)
				continue
			}

			pubkey, ok := key.(*rsa.PublicKey)
			if !ok {
				log.Printf("WARN: %s is not a valid RSA public key", filename)
				continue
			}

			keys[clientID] = pubkey
		}
	}

	return keys, nil
}
