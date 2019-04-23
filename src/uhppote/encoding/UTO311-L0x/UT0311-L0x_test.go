package UTO311_L0x

import (
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"time"
	"uhppote/types"
)

func TestMarshal(t *testing.T) {
	expected := []byte{
		0x17, 0x5f, 0x7d, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd2, 0x04, 0x01, 0x00, 0xc0, 0xa8, 0x01, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92,
		0x20, 0x18, 0x08, 0x16, 0x20, 0x19, 0x04, 0x16, 0x12, 0x34, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	mac, _ := net.ParseMAC("00:66:19:39:55:2d")
	date, _ := time.ParseInLocation("2006-01-02", "2018-08-16", time.Local)
	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-04-16 12:34:56", time.Local)

	request := struct {
		MsgType      types.MsgType      `uhppote:"value:0x5f"`
		Byte         byte               `uhppote:"offset:2"`
		Uint32       uint32             `uhppote:"offset:4"`
		Uint16       uint16             `uhppote:"offset:8"`
		True         bool               `uhppote:"offset:10"`
		False        bool               `uhppote:"offset:11"`
		Address      net.IP             `uhppote:"offset:12"`
		MacAddress   net.HardwareAddr   `uhppote:"offset:20"`
		SerialNumber types.SerialNumber `uhppote:"offset:26"`
		Version      types.Version      `uhppote:"offset:30"`
		Date         types.Date         `uhppote:"offset:32"`
		DateTime     types.DateTime     `uhppote:"offset:36"`
	}{
		Byte:         0x7d,
		Uint32:       423187757,
		Uint16:       1234,
		True:         true,
		False:        false,
		Address:      net.IPv4(192, 168, 1, 2),
		MacAddress:   mac,
		SerialNumber: 423187757,
		Version:      0x0892,
		Date:         types.Date{date},
		DateTime:     types.DateTime{datetime},
	}

	m, err := Marshal(request)

	if err != nil {
		t.Errorf("Marshal returned unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Marshal returned invalid message - \nExpected:\n%s\nReturned:\n%s", print(expected), print(m))
		return
	}
}

func TestUnmarshal(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x6e, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xd2, 0x04, 0x00, 0x00, 0xc0, 0xa8, 0x00, 0x00,
		0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55,
		0x39, 0x19, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType      types.MsgType      `uhppote:"offset:1, value:0x94"`
		Byte         byte               `uhppote:"offset:2"`
		Uint32       uint32             `uhppote:"offset:4"`
		Uint16       uint16             `uhppote:"offset:8"`
		Address      net.IP             `uhppote:"offset:12"`
		SubnetMask   net.IP             `uhppote:"offset:16"`
		Gateway      net.IP             `uhppote:"offset:20"`
		MacAddress   net.HardwareAddr   `uhppote:"offset:24"`
		SerialNumber types.SerialNumber `uhppote:"offset:30"`
		Version      types.Version      `uhppote:"offset:34"`
		Date         types.Date         `uhppote:"offset:36"`
		DateTime     types.DateTime     `uhppote:"offset:40"`
		True         bool               `uhppote:"offset:47"`
		False        bool               `uhppote:"offset:48"`
	}{}

	err := Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Expected 'byte':0x%02X, got: 0x%02X\n", 0x94, reply.MsgType)
	}

	if reply.Byte != 0x6e {
		t.Errorf("Expected 'byte':%02x, got: %02x\n", 0x6e, reply.Byte)
	}

	if reply.Uint32 != 423187757 {
		t.Errorf("Expected 'uint32':%v, got: %v\n", 423187757, reply.Uint32)
	}

	if reply.Uint16 != 1234 {
		t.Errorf("Expected 'uint16':%v, got: %v\n", 1234, reply.Uint16)
	}

	if !reflect.DeepEqual(reply.Address, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Expected IP address '%v', got: '%v'\n", net.IPv4(192, 168, 0, 0), reply.Address)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(255, 255, 255, 0), reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Expected subnet mask '%v', got: '%v'\n", net.IPv4(0, 0, 0, 0), reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, MAC) {
		t.Errorf("Expected MAC address '%v', got: '%v'\n", MAC, reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Expected version '0x%04X', got: '0x%04X'\n", 0x0892, reply.Version)
	}

	date, _ := time.ParseInLocation("2006-01-02", "2018-08-16", time.Local)
	if reply.Date.Date != date {
		t.Errorf("Expected date '%v', got: '%v'\n", date, reply.Date)
	}

	datetime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-12-31 12:23:34", time.Local)
	if reply.DateTime.DateTime != datetime {
		t.Errorf("Expected date '%v', got: '%v'\n", datetime, reply.DateTime)
	}

	if reply.True != true {
		t.Errorf("Expected door 1 '%v', got: '%v\n", true, reply.True)
	}

	if reply.False != false {
		t.Errorf("Expected door 2 '%v', got: '%v\n", false, reply.False)
	}
}

func TestUnmarshalWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x92,
		0x20, 0x18, 0x08, 0x16, 0x20, 0x18, 0x12, 0x31, 0x12, 0x23, 0x34, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := struct {
		MsgType types.MsgType `uhppote:"offset:1, value:0x92"`
		Uint32  uint32        `uhppote:"offset:4"`
	}{}

	err := Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", " Invalid value in message - expected 92, received 0x94")
		return
	}
}

func print(m []byte) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), "$1"))
}
