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

func genkeys(rsa, hotp string) error {
	encryption := filepath.Join(rsa, "mqtt", "rsa", "encryption")
	signing := filepath.Join(rsa, "mqtt", "rsa", "signing")

	// ... encryption keys
	pk := filepath.Join(encryption, "mqttd.key")
	genkey(&pk, nil, "   ... creating RSA encryption key")

	// ... signing keys
	pk = filepath.Join(signing, "mqttd.key")
	genkey(&pk, nil, "   ... creating RSA signing key")

	// ... event message keys
	pubkey := filepath.Join(encryption, "events.pub")
	genkey(nil, &pubkey, "   ... creating RSA event message key")

	// ... system message keys
	pubkey = filepath.Join(encryption, "system.pub")
	genkey(nil, &pubkey, "   ... creating RSA system message key")

	// ... HOTP secrets

	_, err := os.Stat(hotp)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... creating HOTP secrets file '%s'\n", hotp)
		f, err := os.Create(hotp)
		if err != nil {
			return err
		}

		f.Close()
	}

	return nil
}

func genkey(pk, pubkey *string, msg string) error {
	reader := rand.Reader
	bits := 2048

	key, err := rsa.GenerateKey(reader, bits)
	if err != nil {
		return err
	}

	if pk != nil {
		_, err := os.Stat(*pk)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if os.IsNotExist(err) {
			fmt.Printf("%s '%s'\n", msg, *pk)
			if err := storePrivateKey(*pk, key); err != nil {
				return err
			}

			if pubkey != nil {
				fmt.Printf("%s '%s'\n", msg, *pubkey)
				if err := storePublicKey(*pubkey, key.PublicKey); err != nil {
					return err
				}
			}
		}
	} else if pubkey != nil {
		_, err := os.Stat(*pubkey)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if os.IsNotExist(err) {
			fmt.Printf("%s '%s'\n", msg, *pubkey)
			if err := storePublicKey(*pubkey, key.PublicKey); err != nil {
				return err
			}
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
