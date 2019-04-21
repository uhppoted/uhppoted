package uhppote

import (
	"reflect"
	"testing"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
)

func TestMarshalGetStatusRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2D, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := GetStatusRequest{
		SerialNumber: 423187757,
	}

	m, err := codec.Marshal(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid byte array:\nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshalGetStatusResponse(t *testing.T) {
	message := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x39, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x01,
		0xaa, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x04, 0x19, 0x17, 0x00, 0x09, 0x06, 0x01, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x01, 0x01, 0x09, 0x14, 0x37, 0x02, 0x11, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00,
		0x2b, 0x01, 0x00, 0x19, 0x04, 0x20, 0x00, 0x00, 0x93, 0x26, 0x04, 0x88, 0x08, 0x92, 0x00, 0x00,
	}

	reply := GetStatusResponse{}

	err := codec.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if reply.MsgType != 0x20 {
		t.Errorf("Incorrect 'message type' - expected:%02X, got:%02x", 0x32, reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' - expected:%v, got:%v", 423187757, reply.SerialNumber)
	}

	if reply.LastIndex != 57 {
		t.Errorf("Incorrect 'last index' - expected:%v, got:%v", 57, reply.LastIndex)
	}

	if reply.SwipeRecord != 1 {
		t.Errorf("Incorrect 'swipe record' - expected:%v, got:%v", 1, reply.SwipeRecord)
	}

	if reply.Granted {
		t.Errorf("Incorrect 'access granted' - expected:%v, got:%v", false, reply.Granted)
	}

	if reply.Door != 3 {
		t.Errorf("Incorrect 'door' - expected:%v, got:%v", 3, reply.Door)
	}

	if !reply.DoorOpen {
		t.Errorf("Incorrect 'door open' - expected:%v, got:%v", true, reply.DoorOpen)
	}

	if reply.CardNumber != 6154410 {
		t.Errorf("Incorrect 'card number' - expected:%v, got:%v", 6154410, reply.CardNumber)
	}

	swiped, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-04-19 17:00:09", time.Local)
	if reply.SwipeDateTime.DateTime != swiped {
		t.Errorf("Incorrect 'swipe date/time' - expected:%s, got:%s", swiped.Format("2006-01-02 15:04:05"), reply.SwipeDateTime)
	}

	if reply.SwipeReason != 6 {
		t.Errorf("Incorrect 'swipe reason' - expected:%v, got:%v", 6, reply.SwipeReason)
	}

	if !reply.Door1State {
		t.Errorf("Incorrect 'door 1 state' - expected:%v, got:%v", true, reply.Door1State)
	}

	if reply.Door2State {
		t.Errorf("Incorrect 'door 2 state' - expected:%v, got:%v", false, reply.Door2State)
	}

	if !reply.Door3State {
		t.Errorf("Incorrect 'door 3 state' - expected:%v, got:%v", true, reply.Door3State)
	}

	if !reply.Door4State {
		t.Errorf("Incorrect 'door 4 state' - expected:%v, got:%v", true, reply.Door4State)
	}

	if reply.Door1Button {
		t.Errorf("Incorrect 'door 1 button' - expected:%v, got:%v", false, reply.Door1Button)
	}

	if reply.Door2Button {
		t.Errorf("Incorrect 'door 2 button' - expected:%v, got:%v", false, reply.Door2Button)
	}

	if !reply.Door3Button {
		t.Errorf("Incorrect 'door 3 button' - expected:%v, got:%v", true, reply.Door3Button)
	}

	if !reply.Door4Button {
		t.Errorf("Incorrect 'door 4 button' - expected:%v, got:%v", true, reply.Door4Button)
	}

	if reply.SystemState != 9 {
		t.Errorf("Incorrect 'system state' - expected:%v, got:%v", 9, reply.SystemState)
	}

	sysdate, _ := time.ParseInLocation("2006-01-02", "2019-04-20", time.Local)
	if reply.SystemDate.Date != sysdate {
		t.Errorf("Incorrect 'system date' - expected:%s, got:%s", sysdate.Format("2006-01-02"), reply.SystemDate.Date)
	}

	systime, _ := time.ParseInLocation("15:04:05", "14:37:02", time.Local)
	if reply.SystemTime.Time != systime {
		t.Errorf("Incorrect 'system time' - expected:%s, got:%s", systime.Format("15:04:05"), reply.SystemTime.Time)
	}

	if reply.PacketNumber != 17 {
		t.Errorf("Incorrect 'packet number' - expected:%v, got:%v", 17, reply.PacketNumber)
	}

	if reply.Backup != 33 {
		t.Errorf("Incorrect 'backup' - expected:%v, got:%v", 33, reply.Backup)
	}

	if reply.SpecialMessage != 43 {
		t.Errorf("Incorrect 'special message' - expected:%v, got:%v", 43, reply.SpecialMessage)
	}

	if !reply.LowBattery {
		t.Errorf("Incorrect 'low battery' - expected:%v, got:%v", true, reply.LowBattery)
	}

	if reply.FireAlarm {
		t.Errorf("Incorrect 'fire alarm' - expected:%v, got:%v", false, reply.FireAlarm)
	}
}

func TestUnmarshalGetStatusResponseWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92,
		0x20, 0x18, 0x08, 0x16, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := GetStatusResponse{}

	err := codec.Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", "Invalid value in message - expected 0x92, received 0x94")
		return
	}
}
