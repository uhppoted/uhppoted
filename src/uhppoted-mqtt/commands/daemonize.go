package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
)

func hmac() (string, error) {
	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789#!")
	chars := make([]byte, 256)
	err := error(nil)

	copy(chars[0:64], charset)
	copy(chars[64:128], charset)
	copy(chars[128:192], charset)
	copy(chars[192:256], charset)

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i := 0; i < len(bytes); i++ {
		bytes[i] = chars[bytes[i]]
	}

	return string(bytes), err
}

// Ref. https://gist.github.com/sdorra/1c95de8cb80da31610d2ad767cd6f251
func genkeys(root string) error {
	reader := rand.Reader
	bits := 2048

	pk := filepath.Join(root, "mqtt", "rsa", "encryption", "mqttd.key")
	pubkey := filepath.Join(root, "mqtt", "rsa", "encryption", "mqttd.pub")
	_, err := os.Stat(pk)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		key, err := rsa.GenerateKey(reader, bits)
		if err != nil {
			return err
		}

		fmt.Printf("   ... creating RSA encryption key '%s'\n", pk)
		storePrivateKey(pk, key)
		storePublicKey(pubkey, key.PublicKey)
	}

	pk = filepath.Join(root, "mqtt", "rsa", "signing", "mqttd.key")
	pubkey = filepath.Join(root, "mqtt", "rsa", "signing", "mqttd.pub")
	_, err = os.Stat(pk)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		key, err := rsa.GenerateKey(reader, bits)
		if err != nil {
			return err
		}

		fmt.Printf("   ... creating RSA signing key '%s'\n", pk)
		storePrivateKey(pk, key)
		storePublicKey(pubkey, key.PublicKey)
	}

	return nil
}

func storePrivateKey(path string, key *rsa.PrivateKey) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	var pemkey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return pem.Encode(f, pemkey)
}

func storePublicKey(path string, key rsa.PublicKey) error {
	bytes, err := asn1.Marshal(key)
	if err != nil {
		return err
	}

	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	return pem.Encode(f, pemkey)
}
