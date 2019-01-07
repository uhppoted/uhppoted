package messages

import (
	"net"
	"reflect"
	"testing"
)

func TestSearchFrom(t *testing.T) {
	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply, err := NewSearch(message)

	if err != nil {
		t.Errorf("SearchFrom returned error from valid message: %v\n", err)
	}

	if reply == nil {
		t.Errorf("SearchFrom returned nil from valid message: %v\n", reply)
	}

	if reply.SOM != 0x17 {
		t.Errorf("SearchFrom returned incorrect 'device type' from valid message: %02x\n", reply.SOM)
	}

	if reply.MsgType != 0x94 {
		t.Errorf("SearchFrom returned incorrect 'message type' from valid message: %02x\n", reply.MsgType)
	}

	if reply.Device.SerialNumber != 423187757 {
		t.Errorf("SearchFrom returned incorrect 'serial number' from valid message: %v\n", reply.Device.SerialNumber)
	}

	if !reflect.DeepEqual(reply.Device.IpAddress, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("SearchFrom returned incorrect 'IP address' from valid message: %v\n", reply.Device.IpAddress)
	}

	if !reflect.DeepEqual(reply.Device.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("SearchFrom returned incorrect 'subnet mask' from valid message: %v\n", reply.Device.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Device.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("SearchFrom returned incorrect 'gateway' from valid message: %v\n", reply.Device.Gateway)
	}

	if !reflect.DeepEqual(reply.Device.MacAddress, MAC) {
		t.Errorf("SearchFrom returned incorrect 'MAC address' from valid message: %v\n", reply.Device.MacAddress)
	}

	if reply.Device.Version != "0892" {
		t.Errorf("SearchFrom returned incorrect 'version' from valid message: %v\n", reply.Device.Version)
	}

	if reply.Device.Date != "2018-08-16" {
		t.Errorf("SearchFrom returned incorrect 'date' from valid message: %v\n", reply.Device.Date)
	}
}
