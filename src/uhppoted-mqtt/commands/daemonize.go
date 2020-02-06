package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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

	// ... encryption keys
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
		if err := storePrivateKey(pk, key); err != nil {
			return err
		}

		if err := storePublicKey(pubkey, key.PublicKey); err != nil {
			return err
		}
	}

	// ... signing keys
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
		if err := storePrivateKey(pk, key); err != nil {
			return err
		}

		if err := storePublicKey(pubkey, key.PublicKey); err != nil {
			return err
		}
	}

	// ... event message keys
	pk = filepath.Join(root, "mqtt", "rsa", "encryption", "events.key")
	pubkey = filepath.Join(root, "mqtt", "rsa", "encryption", "events.pub")
	_, err = os.Stat(pk)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		key, err := rsa.GenerateKey(reader, bits)
		if err != nil {
			return err
		}

		fmt.Printf("   ... creating RSA event message key '%s'\n", pk)
		if err := storePublicKey(pubkey, key.PublicKey); err != nil {
			return err
		}
	}

	// ... system message keys
	pk = filepath.Join(root, "mqtt", "rsa", "encryption", "system.key")
	pubkey = filepath.Join(root, "mqtt", "rsa", "encryption", "system.pub")
	_, err = os.Stat(pk)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		key, err := rsa.GenerateKey(reader, bits)
		if err != nil {
			return err
		}

		fmt.Printf("   ... creating RSA system message key '%s'\n", pk)
		if err := storePublicKey(pubkey, key.PublicKey); err != nil {
			return err
		}
	}

	return nil
}

func storePrivateKey(path string, key *rsa.PrivateKey) error {
	bytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(path), os.ModePerm)

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	var pemkey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: bytes,
	}

	return pem.Encode(f, pemkey)
}

func storePublicKey(path string, key rsa.PublicKey) error {
	bytes, err := x509.MarshalPKIXPublicKey(&key)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(path), os.ModePerm)

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
