package conf

import (
	"net"
	"reflect"
	"testing"
)

type testType struct {
	value string
}

type device struct {
	address string
}

func (t *testType) UnmarshalConf(s string) (interface{}, error) {
	return &testType{s}, nil
}

var configuration = []byte(`# key-value pairs
udp.address = 192.168.1.100:54321
interface.value = qwerty
interface.pointer = uiop
sys.enabled = true
sys.integer = -13579
sys.unsigned = 8081
sys.string = asdfghjkl

# DEVICES
UT0311-L0x.305419896.address = 192.168.1.100:60000
UT0311-L0x.305419896.door.1 = Front Door
UT0311-L0x.305419896.door.2 = Side Door
UT0311-L0x.305419896.door.3 = Garage
UT0311-L0x.305419896.door.4 = Workshop
`)

func TestUnmarshal(t *testing.T) {
	config := struct {
		UdpAddress *net.UDPAddr    `conf:"udp.address"`
		Interface  testType        `conf:"interface.value"`
		InterfaceP *testType       `conf:"interface.pointer"`
		Enabled    bool            `conf:"sys.enabled"`
		Integer    int             `conf:"sys.integer"`
		Unsigned   uint            `conf:"sys.unsigned"`
		String     string          `conf:"sys.string"`
		Devices    map[uint]device `conf:"UT0311-L0x.*"`
	}{}

	err := Unmarshal(configuration, &config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	address, _ := net.ResolveUDPAddr("udp", "192.168.1.100:54321")
	if !reflect.DeepEqual(config.UdpAddress, address) {
		t.Errorf("Expected 'udp.address' %s, got: %s", address, config.UdpAddress)
	}

	if config.Interface.value != "qwerty" {
		t.Errorf("Expected 'interface' value '%s', got: '%v'", "qwerty", config.Interface)
	}

	if config.InterfaceP == nil || config.InterfaceP.value != "uiop" {
		t.Errorf("Expected 'interface pointer' value '%s', got: '%v'", "uiop", config.InterfaceP)
	}

	if !config.Enabled {
		t.Errorf("Expected 'boolean' value '%v', got: '%v'", true, config.Enabled)
	}

	if config.Integer != -13579 {
		t.Errorf("Expected 'integer' value '%v', got: '%v'", -13579, config.Integer)
	}

	if config.Unsigned != 8081 {
		t.Errorf("Expected 'unsigned' value '%v', got: '%v'", 8081, config.Unsigned)
	}

	if config.String != "asdfghjkl" {
		t.Errorf("Expected 'string' value '%v', got: '%v'", "asdfghjkl", config.String)
	}

	//if _, ok := config.Devices[305419896]; !ok {
	//	t.Errorf("Expected 'device' for ID '%v', got: '%v'", 305419896, false)
	//}
}
