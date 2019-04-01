package messages

import "testing"

func TestNewGetTime(t *testing.T) {
	message := []byte{
		0x17, 0x32, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x20, 0x19, 0x01, 0x02, 0x06, 0x04, 0x20, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply, err := NewGetTime(message)

	if err != nil {
		t.Errorf("NewGetTime returned error from valid message: %v\n", err)
	}

	if reply == nil {
		t.Errorf("NewGetTime returned nil from valid message: %v\n", reply)
	}

	if reply.StartOfMessage != 0x17 {
		t.Errorf("NewGetTime returned incorrect 'start of message' from valid message: %02x\n", reply.StartOfMessage)
	}

	if reply.MsgType != 0x32 {
		t.Errorf("NewGetTime returned incorrect 'message type' from valid message: %02x\n", reply.MsgType)
	}

	// if reply.DateTime.SerialNumber != 423187757 {
	//	t.Errorf("NewGetTime returned incorrect 'serial number' from valid message: %v\n", reply.DateTime.SerialNumber)
	//}

	if reply.DateTime.DateTime.Format("2006-01-02 15:04:05") != "2019-01-02 06:04:20" {
		t.Errorf("NewGetTime returned incorrect 'date/time' from valid message: %v\n", reply.DateTime.DateTime)
	}
}
