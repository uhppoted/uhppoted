package messages

import (
	"reflect"
	"testing"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

func TestMarshalGetStatusRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
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
		t.Errorf("Invalid byte array:\nExpected:\n%s\nReturned:\n%s", dump(expected, ""), dump(m, ""))
		return
	}
}

func TestFactoryUnmarshalGetStatusRequest(t *testing.T) {
	message := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request, err := UnmarshalRequest(message)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if request == nil {
		t.Fatalf("Unexpected request: %v\n", request)
	}

	rq, ok := request.(*GetStatusRequest)
	if !ok {
		t.Fatalf("Invalid request type - expected:%T, got: %T\n", &GetStatusRequest{}, request)
	}

	if rq.MsgType != 0x20 {
		t.Errorf("Incorrect 'message type' from valid message: %02x\n", rq.MsgType)
	}
}

func TestUnmarshalGetStatusResponse(t *testing.T) {
	message := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x39, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x01,
		0xaa, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x04, 0x19, 0x17, 0x00, 0x09, 0x06, 0x01, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x01, 0x01, 0x09, 0x14, 0x37, 0x02, 0x11, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00,
		0x2b, 0x04, 0x01, 0x19, 0x04, 0x20, 0x00, 0x00, 0x93, 0x26, 0x04, 0x88, 0x08, 0x92, 0x00, 0x00,
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

	if reply.EventType != 1 {
		t.Errorf("Incorrect 'event type' - expected:%v, got:%v", 1, reply.EventType)
	}

	if reply.Granted {
		t.Errorf("Incorrect 'access granted' - expected:%v, got:%v", false, reply.Granted)
	}

	if reply.Door != 3 {
		t.Errorf("Incorrect 'door' - expected:%v, got:%v", 3, reply.Door)
	}

	if !reply.DoorOpened {
		t.Errorf("Incorrect 'door opened' - expected:%v, got:%v", true, reply.DoorOpened)
	}

	if reply.UserID != 6154410 {
		t.Errorf("Incorrect 'user ID' - expected:%v, got:%v", 6154410, reply.UserID)
	}

	swiped, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-04-19 17:00:09", time.Local)
	if reply.EventTimestamp != types.DateTime(swiped) {
		t.Errorf("Incorrect 'event timestamp' - expected:%s, got:%s", swiped.Format("2006-01-02 15:04:05"), reply.EventTimestamp)
	}

	if reply.EventResult != 6 {
		t.Errorf("Incorrect 'event result' - expected:%v, got:%v", 6, reply.EventResult)
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
	if reply.SystemDate != types.SystemDate(sysdate) {
		t.Errorf("Incorrect 'system date' - expected:%s, got:%s", sysdate.Format("2006-01-02"), reply.SystemDate.String())
	}

	systime, _ := time.ParseInLocation("15:04:05", "14:37:02", time.Local)
	if reply.SystemTime != types.SystemTime(systime) {
		t.Errorf("Incorrect 'system time' - expected:%s, got:%s", systime.Format("15:04:05"), reply.SystemTime.String())
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

	if reply.Battery != 0x04 {
		t.Errorf("Incorrect 'battery status' - expected:%v, got:%v", 0x04, reply.Battery)
	}

	if reply.FireAlarm != 0x01 {
		t.Errorf("Incorrect 'fire alarm' - expected:%v, got:%v", 0x01, reply.FireAlarm)
	}
}

func TestFactoryUnmarshalGetStatusResponse(t *testing.T) {
	message := []byte{
		0x17, 0x20, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x39, 0x00, 0x00, 0x00, 0x01, 0x00, 0x03, 0x01,
		0xaa, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x04, 0x19, 0x17, 0x00, 0x09, 0x06, 0x01, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x01, 0x01, 0x09, 0x14, 0x37, 0x02, 0x11, 0x00, 0x00, 0x00, 0x21, 0x00, 0x00, 0x00,
		0x2b, 0x04, 0x01, 0x19, 0x04, 0x20, 0x00, 0x00, 0x93, 0x26, 0x04, 0x88, 0x08, 0x92, 0x00, 0x00,
	}

	response, err := UnmarshalResponse(message)

	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	if response == nil {
		t.Fatalf("Unexpected response: %v\n", response)
	}

	reply, ok := response.(*GetStatusResponse)
	if !ok {
		t.Fatalf("Invalid response type - expected:%T, got: %T\n", &GetStatusResponse{}, response)
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

	if reply.EventType != 1 {
		t.Errorf("Incorrect 'event type' - expected:%v, got:%v", 1, reply.EventType)
	}

	if reply.Granted {
		t.Errorf("Incorrect 'access granted' - expected:%v, got:%v", false, reply.Granted)
	}

	if reply.Door != 3 {
		t.Errorf("Incorrect 'door' - expected:%v, got:%v", 3, reply.Door)
	}

	if !reply.DoorOpened {
		t.Errorf("Incorrect 'door opened' - expected:%v, got:%v", true, reply.DoorOpened)
	}

	if reply.UserID != 6154410 {
		t.Errorf("Incorrect 'user ID' - expected:%v, got:%v", 6154410, reply.UserID)
	}

	swiped, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-04-19 17:00:09", time.Local)
	if reply.EventTimestamp != types.DateTime(swiped) {
		t.Errorf("Incorrect 'event timestamp' - expected:%s, got:%s", swiped.Format("2006-01-02 15:04:05"), reply.EventTimestamp)
	}

	if reply.EventResult != 6 {
		t.Errorf("Incorrect 'event result' - expected:%v, got:%v", 6, reply.EventResult)
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
	if reply.SystemDate != types.SystemDate(sysdate) {
		t.Errorf("Incorrect 'system date' - expected:%s, got:%s", sysdate.Format("2006-01-02"), reply.SystemDate.String())
	}

	systime, _ := time.ParseInLocation("15:04:05", "14:37:02", time.Local)
	if reply.SystemTime != types.SystemTime(systime) {
		t.Errorf("Incorrect 'system time' - expected:%s, got:%s", systime.Format("15:04:05"), reply.SystemTime.String())
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

	if reply.Battery != 0x04 {
		t.Errorf("Incorrect 'battery status' - expected:%v, got:%v", 0x04, reply.Battery)
	}

	if reply.FireAlarm != 0x01 {
		t.Errorf("Incorrect 'fire alarm' - expected:%v, got:%v", 0x01, reply.FireAlarm)
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
