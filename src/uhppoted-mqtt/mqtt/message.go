package mqtt

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type msgType int

const (
	msgReply msgType = iota + 1
	msgError
	msgEvent
)

func (mqttd *MQTTD) wrap(msgtype msgType, content interface{}, destID *string) ([]byte, error) {
	bytes, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	signature, err := mqttd.sign(bytes)
	if err != nil {
		return nil, err
	}

	ciphertext, iv, key, err := mqttd.encrypt(bytes, destID)
	if err != nil {
		return nil, err
	}

	message := struct {
		Signature string `json:"signature,omitempty"`
		Key       string `json:"key,omitempty"`
		Reply json.RawMessage `json:"reply,omitempty"`
		Error json.RawMessage `json:"error,omitempty"`
		Event json.RawMessage `json:"event,omitempty"`
	}{
		Signature: base64.StdEncoding.EncodeToString(signature),
		Key:       base64.StdEncoding.EncodeToString(key),
	}

	crypttext, err := json.Marshal(base64.StdEncoding.EncodeToString(append(iv, ciphertext...)))
	if err != nil {
		return nil, err
	}

	switch msgtype {
	case msgReply:
		message.Reply = crypttext
	case msgError:
		message.Error = crypttext
	case msgEvent:
		message.Event = crypttext
	}

	bytes, err = json.Marshal(message)
	if err != nil {
		return nil, err
	}

	payload := struct {
		Message json.RawMessage `json:"message"`
		HMAC    string          `json:"hmac,omitempty"`
	}{
		Message: bytes,
		HMAC:    hex.EncodeToString(mqttd.HMAC.MAC(bytes)),
	}

	bytes, err = json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (mqttd *MQTTD) unwrap(payload []byte) ([]byte, error) {
	message := struct {
		Message json.RawMessage `json:"message"`
		HMAC    *string         `json:"hmac"`
	}{}

	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, fmt.Errorf("Error unmarshaling message (%v)", err)
	}

	if err := mqttd.verify(message.Message, message.HMAC); err != nil {
		return nil, fmt.Errorf("Invalid message (%v)", err)
	}

	body := struct {
		Signature *string         `json:"signature"`
		Key       *string         `json:"key"`
		IV        string          `json:"iv"`
		Request   json.RawMessage `json:"request"`
	}{}

	if err := json.Unmarshal(message.Message, &body); err != nil {
		return nil, fmt.Errorf("Error unmarshaling message body (%v)", err)
	}

	request := []byte(body.Request)

	if body.Key != nil && isBase64(body.Request) {
		plaintext, err := mqttd.decrypt(request, body.IV, *body.Key)
		if err != nil || plaintext == nil {
			return nil, fmt.Errorf("Error decrypting message (%v::%v)", err, plaintext)
		}

		request = plaintext
	}

	misc := struct {
		ClientID  *string `json:"client-id"`
		RequestID *string `json:"request-id"`
		ReplyTo   *string `json:"reply-to"`
		Nonce     *uint64 `json:"nonce"`
	}{}

	if err := json.Unmarshal(request, &misc); err != nil {
		return nil, fmt.Errorf("Error unmarshaling request meta-info (%v)", err)
	}

	authenticated, err := mqttd.authenticate(misc.ClientID, request, body.Signature)
	if err != nil {
		return nil, fmt.Errorf("Error authenticating request (%v)", err)
	}

	if authenticated {
		if err := mqttd.Nonce.Validate(misc.ClientID, misc.Nonce); err != nil {
			return nil, fmt.Errorf("Error validating nonce (%v)", err)
		}
	}

	return request, nil
}

func (m *MQTTD) verify(message []byte, mac *string) error {
	if m.HMAC.Required && mac == nil {
		return errors.New("HMAC required but not present")
	}

	if mac != nil {
		hmac, err := hex.DecodeString(*mac)
		if err != nil {
			return err
		}

		if !m.HMAC.Verify(message, hmac) {
			return errors.New("incorrect HMAC")
		}
	}

	return nil
}

func (m *MQTTD) decrypt(request []byte, iv string, key string) ([]byte, error) {
	var crypttext string = ""

	err := json.Unmarshal(request, &crypttext)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(crypttext)
	if err != nil {
		return nil, fmt.Errorf("Invalid ciphertext (%v)", err)
	}

	keyv, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(key, " ", ""))
	if err != nil {
		return nil, fmt.Errorf("Invalid key (%v)", err)
	}

	ivv, err := hex.DecodeString(iv)
	if err != nil {
		return nil, fmt.Errorf("Invalid IV (%v)", err)
	}

	return m.RSA.Decrypt(append(ivv, ciphertext...), keyv)
}

func (m *MQTTD) authenticate(clientID *string, request []byte, signature *string) (bool, error) {
	if m.Authentication == "HOTP" {
		rq := struct {
			HOTP *string `json:"hotp"`
		}{}

		if clientID == nil {
			return false, errors.New("Invalid request: missing client-id")
		}

		if err := json.Unmarshal(request, &rq); err != nil {
			return false, err
		}

		if rq.HOTP == nil {
			return false, errors.New("Invalid request: missing HOTP token")
		}

		if err := m.HOTP.Validate(*clientID, *rq.HOTP); err != nil {
			return false, err
		}

		return true, nil
	}

	if m.Authentication == "RSA" {
		if clientID == nil {
			return false, errors.New("Invalid request: missing client-id")
		}

		if signature == nil {
			return false, errors.New("Invalid request: missing RSA signature")
		}

		s, err := base64.StdEncoding.DecodeString(*signature)
		if err != nil {
			return false, fmt.Errorf("Invalid request: undecodable RSA signature (%v)", err)
		}

		if err := m.RSA.Validate(*clientID, request, s); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (m *MQTTD) sign(reply []byte) ([]byte, error) {
	if m.SignOutgoing {
		return m.RSA.Sign(reply)
	}

	return nil, nil
}

func (m *MQTTD) encrypt(plaintext []byte, clientID *string) ([]byte, []byte, []byte, error) {
	if m.EncryptOutgoing {
		if clientID == nil {
			return nil, nil, nil, fmt.Errorf("Missing client ID")
		}

		ciphertext, key, err := m.RSA.Encrypt(plaintext, *clientID)
		if err != nil {
			return nil, nil, nil, err
		}

		return ciphertext[16:], ciphertext[:16], key, nil
	}

	return plaintext, nil, nil, nil
}
