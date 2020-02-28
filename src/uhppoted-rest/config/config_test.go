package config

import (
	"github.com/uhppoted/uhppote-core/encoding/conf"
	"net"
	"reflect"
	"testing"
)

var configuration = []byte(`# UDP
bind.address = 192.168.1.100:54321
broadcast.address = 192.168.1.255:60001

# REST API              
rest.CORS.enabled = true
rest.http.enabled = true
rest.http.port = 8081
rest.https.enabled = false
rest.https.port = 8883
rest.tls.key = qwerty.key
rest.tls.certificate = qwerty.cert
rest.tls.ca = qwerty.ca

# OPEN API
openapi.enabled = true
openapi.directory = asdfghjkl

# DEVICES
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Front Door
UT0311-L0x.405419896.door.2 = Side Door
UT0311-L0x.405419896.door.3 = Garage
UT0311-L0x.405419896.door.4 = Workshop
`)

func TestUnmarshal(t *testing.T) {
	config := NewConfig()

	err := conf.Unmarshal(configuration, config)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	address := net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 54321,
		Zone: "",
	}

	if !reflect.DeepEqual(config.BindAddress, &address) {
		t.Errorf("Expected 'bind.address' %s, got: %s", &address, config.BindAddress)
	}

	address = net.UDPAddr{
		IP:   []byte{192, 168, 1, 255},
		Port: 60001,
		Zone: "",
	}

	if !reflect.DeepEqual(config.BroadcastAddress, &address) {
		t.Errorf("Expected 'broadcast.address' %s, got:%s", &address, config.BroadcastAddress)
	}

	if !config.HttpEnabled {
		t.Errorf("Expected 'REST::HttpEnabled' %v, got:%v", true, config.HttpEnabled)
	}

	if config.HttpPort != 8081 {
		t.Errorf("Expected 'REST::HttpPort' %v, got:%v", 8081, config.HttpPort)
	}

	if config.HttpsEnabled {
		t.Errorf("Expected 'REST::HttpsEnabled' %v, got:%v", false, config.HttpsEnabled)
	}

	if config.HttpsPort != 8883 {
		t.Errorf("Expected 'REST::HttpsPort' %v, got:%v", 8883, config.HttpsPort)
	}

	if config.TLSKeyFile != "qwerty.key" {
		t.Errorf("Expected 'REST::TLSKeyFile' %v, got:%v", "qwerty.key", config.TLSKeyFile)
	}

	if config.TLSCertificateFile != "qwerty.cert" {
		t.Errorf("Expected 'REST::TLSCertificateFile' %v, got:%v", "qwerty.cert", config.TLSCertificateFile)
	}

	if config.CACertificateFile != "qwerty.ca" {
		t.Errorf("Expected 'REST::CACertificateFile' %v, got:%v", "qwerty.ca", config.TLSCertificateFile)
	}

	if !config.CORSEnabled {
		t.Errorf("Expected 'REST::CORSEnabled' %v, got:%v", true, config.CORSEnabled)
	}

	if !config.OpenApi.Enabled {
		t.Errorf("Expected 'OpenApi::Enabled' %v, got:%v", true, config.OpenApi.Enabled)
	}

	if config.OpenApi.Directory != "asdfghjkl" {
		t.Errorf("Expected 'OpenApi::Directory' %v, got:%v", "asdfghjkl", config.OpenApi.Directory)
	}

	if d, _ := config.Devices[405419896]; d == nil {
		t.Errorf("Expected 'device' for ID '%v', got:'%v'", 405419896, d)
	} else {
		address = net.UDPAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 60000,
			Zone: "",
		}

		if !reflect.DeepEqual(d.Address, &address) {
			t.Errorf("Expected 'device.address' %s for ID '%v', got:'%v'", &address, 405419896, d.Address)
		}

		if len(d.Door) != 4 {
			t.Errorf("Expected 4 entries for 'device.door' %s for ID '%v', got:%v", &address, 405419896, len(d.Door))
		} else {
			if d.Door[0] != "Front Door" {
				t.Errorf("Expected 'device.door[0]' %s for ID '%v', got:'%s'", "Front Door", 405419896, d.Door[0])
			}

			if d.Door[1] != "Side Door" {
				t.Errorf("Expected 'device.door[1]' %s for ID '%v', got:'%s'", "Side Door", 405419896, d.Door[1])
			}

			if d.Door[2] != "Garage" {
				t.Errorf("Expected 'device.door[2]' %s for ID '%v', got:'%s'", "Garage", 405419896, d.Door[2])
			}

			if d.Door[3] != "Workshop" {
				t.Errorf("Expected 'device.door[3]' %s for ID '%v', got:'%s'", "Workshop", 405419896, d.Door[3])
			}
		}
	}
}
