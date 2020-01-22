package config

import (
	"net"
	"reflect"
	"testing"
	"uhppote/encoding/conf"
)

var configuration = []byte(`# UDP
bind.address = 192.168.1.100:54321
broadcast.address = 192.168.1.255:60001
listen.address = 192.168.1.100:12345

# MQTT
mqtt.broker = tls://127.0.0.63:8887
mqtt.topic.root = twystd-qwerty
mqtt.topic.replies = /uiop
mqtt.topic.events = ./asdf
mqtt.topic.system = sys
mqtt.broker.certificate = mqtt-broker.cert
mqtt.client.certificate = mqtt-client.cert
mqtt.client.key = mqtt-client.key

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

	address = net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 12345,
		Zone: "",
	}

	if !reflect.DeepEqual(config.ListenAddress, &address) {
		t.Errorf("Expected 'listen.address' %s, got: %s", &address, config.ListenAddress)
	}

	if config.Broker != "tls://127.0.0.63:8887" {
		t.Errorf("Expected 'mqtt.broker' %s, got:%v", "tls://127.0.0.63:8887", config.Broker)
	}

	if config.BrokerCertificate != "mqtt-broker.cert" {
		t.Errorf("Expected 'mqtt.broker.certificate' %s, got:%v", "mqtt-broker.cert", config.BrokerCertificate)
	}

	if config.ClientCertificate != "mqtt-client.cert" {
		t.Errorf("Expected 'mqtt.client.certificate' %s, got:%v", "mqtt-client.cert", config.ClientCertificate)
	}

	if config.ClientKey != "mqtt-client.key" {
		t.Errorf("Expected 'mqtt.client.key' %s, got:%v", "mqtt-client.key", config.ClientKey)
	}

	if config.Topics.Root != "twystd-qwerty" {
		t.Errorf("Expected 'mqtt::topic' %v, got:%v", "twystd-qwerty", config.Topics.Root)
	}

	if config.Topics.Resolve(config.Topics.Requests) != "twystd-qwerty/requests" {
		t.Errorf("Expected 'mqtt::topic.requests' %v, got:%v", "wystd-qwerty/requests", config.Topics.Resolve(config.Topics.Requests))
	}

	if config.Topics.Resolve(config.Topics.Replies) != "uiop" {
		t.Errorf("Expected 'mqtt::topic.replies' %v, got:%v", "uiop", config.Topics.Resolve(config.Topics.Replies))
	}

	if config.Topics.Resolve(config.Topics.Events) != "twystd-qwerty/asdf" {
		t.Errorf("Expected 'mqtt::topic.events' %v, got:%v", "twystd-qwerty/asdf", config.Topics.Resolve(config.Topics.Events))
	}

	if config.Topics.Resolve(config.Topics.System) != "twystd-qwerty/sys" {
		t.Errorf("Expected 'mqtt::topic.system' %v, got:%v", "twystd-qwerty/sys", config.Topics.Resolve(config.Topics.System))
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
