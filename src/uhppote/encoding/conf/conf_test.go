package conf

import (
	"net"
	"reflect"
	"testing"
)

type testType struct {
	value string
}

func (t *testType) UnmarshalConf(s string) (interface{}, error) {
	return &testType{s}, nil
}

var configuration = []byte(`# UDP
udp.address = 192.168.1.100:54321
interface.value = qwerty
interface.pointer = uiop

# REST API
rest.enabled = false
rest.port = 8080
rest.certificate = /etc/uhppoted/rest/uhppoted.cert

# DEVICES
UT0311-L0x.305419896.address = 192.168.1.100:60000
UT0311-L0x.305419896.door.1 = Front Door
UT0311-L0x.305419896.door.2 = Side Door
UT0311-L0x.305419896.door.3 = Garage
UT0311-L0x.305419896.door.4 = Workshop
`)

func TestUnmarshal(t *testing.T) {
	config := struct {
		UdpAddress *net.UDPAddr `conf:"udp.address"`
		Interface  testType     `conf:"interface.value"`
		InterfaceP *testType    `conf:"interface.pointer"`
	}{}

	err := Unmarshal(configuration, &config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	address, _ := net.ResolveUDPAddr("udp", "192.168.1.100:54321")
	if !reflect.DeepEqual(config.UdpAddress, address) {
		t.Errorf("Expected 'udp.address' %s, got: %s\n", address, config.UdpAddress)
	}

	if config.Interface.value != "qwerty" {
		t.Errorf("Expected interface value '%s', got: '%v'\n", "qwerty", config.Interface)
	}

	if config.InterfaceP == nil || config.InterfaceP.value != "uiop" {
		t.Errorf("Expected interface pointer value '%s', got: '%v'\n", "uiop", config.InterfaceP)
	}
}
