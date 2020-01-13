package mqtt

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func unwrap(mqttd *MQTTD, payload []byte) ([]byte, error) {
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

	return message.Message, nil
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

		ciphertext, iv, key, err := m.RSA.Encrypt(plaintext, *clientID)
		if err != nil {
			return nil, nil, nil, err
		}

		crypttext, err := json.Marshal(base64.StdEncoding.EncodeToString(ciphertext))
		if err != nil {
			return nil, nil, nil, err
		}

		return crypttext, iv, key, nil
	}

	return plaintext, nil, nil, nil
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

	return m.RSA.Decrypt(ciphertext, ivv, keyv)
}

func (m *MQTTD) authenticate(clientID *string, request []byte, signature *string) error {
	if m.Authentication == "HOTP" {
		rq := struct {
			HOTP *string `json:"hotp"`
		}{}

		if clientID == nil {
			return errors.New("Invalid request: missing client-id")
		}

		if err := json.Unmarshal(request, &rq); err != nil {
			return err
		}

		if rq.HOTP == nil {
			return errors.New("Invalid request: missing HOTP token")
		}

		return m.HOTP.Validate(*clientID, *rq.HOTP)
	}

	if m.Authentication == "RSA" {
		rq := struct {
			SequenceNo *uint64 `json:"sequence-no"`
		}{}

		if clientID == nil {
			return errors.New("Invalid request: missing client-id")
		}

		if signature == nil {
			return errors.New("Invalid request: missing RSA signature")
		}

		s, err := base64.StdEncoding.DecodeString(*signature)
		if err != nil {
			return fmt.Errorf("Invalid request: undecodable RSA signature (%v)", err)
		}

		if err := json.Unmarshal(request, &rq); err != nil {
			return err
		}

		if rq.SequenceNo == nil {
			return errors.New("Invalid request: missing sequence number")
		}

		return m.RSA.Validate(*clientID, request, s, *rq.SequenceNo)
	}

	return nil
}

// func (m *MQTTD) reply(ctx context.Context, response interface{}) {
// 	client, ok := ctx.Value("client").(MQTT.Client)
// 	if !ok {
// 		panic("MQTT client not included in context")
// 	}
//
// 	topic, ok := ctx.Value("topic").(string)
// 	if !ok {
// 		panic("MQTT root topic not included in context")
// 	}
//
// 	replyTo := "reply"
// 	if rq, ok := ctx.Value("request").(request); ok {
// 		if rq.ReplyTo != nil {
// 			replyTo = *rq.ReplyTo
// 		}
// 	}
//
// 	r, err := json.Marshal(response)
// 	if err != nil {
// 		ctx.Value("log").(*log.Logger).Printf("WARN  %v", err)
// 		return
// 	}
//
// 	signature, err := m.RSA.Sign(r)
// 	if err != nil {
// 		ctx.Value("log").(*log.Logger).Printf("WARN  %v", err)
// 		return
// 	}
//
// 	msg := struct {
// 		ServerID  string          `json:"server-id,omitempty"`
// 		Signature string          `json:"signature,omitempty"`
// 		Key       string          `json:"key,omitempty"`
// 		IV        string          `json:"iv,omitempty"`
// 		Reply     json.RawMessage `json:"reply,omitempty"`
// 	}{
// 		ServerID:  "twystd-uhppoted",
// 		Signature: hex.EncodeToString(signature),
// 		Reply:     r,
// 	}
//
// 	msgbytes, err := json.Marshal(msg)
// 	if err != nil {
// 		ctx.Value("log").(*log.Logger).Printf("WARN  %v", err)
// 		return
// 	}
//
// 	hmac := hex.EncodeToString(m.HMAC.MAC(msgbytes))
//
// 	message := struct {
// 		Message json.RawMessage `json:"message"`
// 		HMAC    string          `json:"hmac,omitempty"`
// 	}{
// 		Message: msgbytes,
// 		HMAC:    hmac,
// 	}
//
// 	b, err := json.Marshal(message)
// 	if err != nil {
// 		ctx.Value("log").(*log.Logger).Printf("WARN  %v", err)
// 		return
// 	}
//
// 	token := client.Publish(topic+"/"+replyTo, 0, false, string(b))
// 	token.Wait()
// }
