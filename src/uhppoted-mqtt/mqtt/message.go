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
	msgSystem
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

	body, key, err := mqttd.encrypt(bytes, destID)
	if err != nil {
		return nil, err
	}

	message := struct {
		Signature string          `json:"signature,omitempty"`
		Key       string          `json:"key,omitempty"`
		Reply     json.RawMessage `json:"reply,omitempty"`
		Error     json.RawMessage `json:"error,omitempty"`
		Event     json.RawMessage `json:"event,omitempty"`
		System    json.RawMessage `json:"system,omitempty"`
	}{
		Signature: base64.StdEncoding.EncodeToString(signature),
		Key:       base64.StdEncoding.EncodeToString(key),
	}

	switch msgtype {
	case msgReply:
		message.Reply = body
	case msgError:
		message.Error = body
	case msgEvent:
		message.Event = body
	case msgSystem:
		message.System = body
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

func (mqttd *MQTTD) unwrap(payload []byte) (*request, error) {
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

	bytes := []byte(body.Request)

	if body.Key != nil && isBase64(body.Request) {
		plaintext, err := mqttd.decrypt(bytes, body.IV, *body.Key)
		if err != nil || plaintext == nil {
			return nil, fmt.Errorf("Error decrypting message (%v::%v)", err, plaintext)
		}

		bytes = plaintext
	}

	misc := struct {
		ClientID  *string `json:"client-id"`
		RequestID *string `json:"request-id"`
		ReplyTo   *string `json:"reply-to"`
		Nonce     *uint64 `json:"nonce"`
	}{}

	if err := json.Unmarshal(bytes, &misc); err != nil {
		return nil, fmt.Errorf("Error unmarshaling request meta-info (%v)", err)
	}

	authenticated, err := mqttd.authenticate(misc.ClientID, bytes, body.Signature)
	if err != nil {
		return nil, err
	}

	if authenticated {
		if err := mqttd.Encryption.Nonce.Validate(misc.ClientID, misc.Nonce); err != nil {
			return nil, fmt.Errorf("Message cannot be authenticated (%v)", err)
		}
	}

	return &request{
		ClientID:  misc.ClientID,
		RequestID: misc.RequestID,
		ReplyTo:   misc.ReplyTo,
		Request:   bytes,
	}, nil
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

func (m *MQTTD) encrypt(plaintext []byte, clientID *string) ([]byte, []byte, error) {
	if m.Encryption.EncryptOutgoing {
		if clientID == nil {
			return nil, nil, fmt.Errorf("Missing client ID")
		}

		ciphertext, key, err := m.Encryption.RSA.Encrypt(plaintext, *clientID, "request")
		if err != nil {
			return nil, nil, err
		}

		crypttext, err := json.Marshal(base64.StdEncoding.EncodeToString(ciphertext))
		if err != nil {
			return nil, nil, err
		}

		return crypttext, key, nil
	}

	return plaintext, nil, nil
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

	return m.Encryption.RSA.Decrypt(append(ivv, ciphertext...), keyv, "request")
}

func (m *MQTTD) authenticate(clientID *string, request []byte, signature *string) (bool, error) {
	if (strings.Contains(m.Authentication, "ANY") || strings.Contains(m.Authentication, "RSA")) && clientID != nil && signature != nil {
		s, err := base64.StdEncoding.DecodeString(*signature)
		if err != nil {
			return false, fmt.Errorf("Invalid request: undecodable RSA signature (%v)", err)
		}

		if err := m.Encryption.RSA.Validate(*clientID, request, s); err != nil {
			return false, err
		}

		return true, nil
	}

	if (strings.Contains(m.Authentication, "ANY") || strings.Contains(m.Authentication, "HOTP")) && clientID != nil {
		rq := struct {
			HOTP *string `json:"hotp"`
		}{}

		if err := json.Unmarshal(request, &rq); err == nil && rq.HOTP != nil {
			if err := m.Encryption.HOTP.Validate(*clientID, *rq.HOTP); err != nil {
				return false, err
			}

			return true, nil
		}
	}
	if strings.Contains(m.Authentication, "NONE") {
		return false, nil
	}

	return false, fmt.Errorf("Could not authenticate %s", *clientID)
}

func (m *MQTTD) sign(reply []byte) ([]byte, error) {
	if m.Encryption.SignOutgoing {
		return m.Encryption.RSA.Sign(reply)
	}

	return nil, nil
}
