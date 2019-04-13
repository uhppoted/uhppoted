package uhppote

import (
	"reflect"
	"testing"
	"time"
	"uhppote/encoding"
)

func TestMarshalGetSwipeRequest(t *testing.T) {
	expected := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := GetSwipeRequest{
		MsgType:      0xb0,
		SerialNumber: 423187757,
		Index:        1,
	}

	m, err := uhppote.Marshal(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid byte array:\nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshalGetSwipeResponse(t *testing.T) {
	message := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0x01,
		0xad, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x02, 0x10, 0x07, 0x12, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0x26, 0x80, 0x39, 0x08, 0x92, 0x00, 0x00,
	}

	reply := GetSwipeResponse{}

	err := uhppote.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if reply.MsgType != 0xb0 {
		t.Errorf("Incorrect 'message type' - expected:%02X, got:%02x\n", 0xb0, reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' - expected:%d, got: %v\n", 423187757, reply.SerialNumber)
	}

	if reply.Index != 8 {
		t.Errorf("Incorrect 'index' - expected:%d, got:%d\n", 8, reply.Index)
	}

	if reply.Type != 2 {
		t.Errorf("Incorrect 'type' - expected:%d, got:%d\n", 2, reply.Type)
	}

	if reply.Granted != true {
		t.Errorf("Incorrect 'granted' - expected:%v, got:%v\n", true, reply.Granted)
	}

	if reply.Door != 3 {
		t.Errorf("Incorrect 'door' - expected:%d, got:%d\n", 3, reply.Door)
	}

	if reply.DoorState != 1 {
		t.Errorf("Incorrect 'door state' - expected:%d, got:%d\n", 1, reply.DoorState)
	}

	if reply.CardNumber != 6154413 {
		t.Errorf("Incorrect 'card number' - expected:%d, got: %v\n", 6154413, reply.CardNumber)
	}

	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-02-10 07:12:01", time.Local)
	if reply.Timestamp.DateTime != timestamp {
		t.Errorf("Incorrect 'timestamp' - expected:%s, got:%s\n", timestamp.Format("2006-01-02 15:04:05"), reply.Timestamp)
	}
}