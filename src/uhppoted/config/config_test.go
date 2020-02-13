package config

import (
	"net"
	"reflect"
	"testing"
	"time"
	"uhppote/encoding/conf"
)

var configuration = []byte(`# SYSTEM
bind.address = 192.168.1.100:54321
broadcast.address = 192.168.1.255:30000
listen.address = 192.168.1.100:12345

monitoring.healthcheck.interval = 31s
monitoring.healthcheck.idle = 67s
monitoring.healthcheck.ignore = 97s
monitoring.watchdog.interval = 23s

# MQTT
mqtt.connection.broker = tls://127.0.0.63:8887
mqtt.connection.client.ID = muppet
mqtt.connection.broker.certificate = mqtt-broker.cert
mqtt.connection.client.certificate = mqtt-client.cert
mqtt.connection.client.key = mqtt-client.key
mqtt.topic.root = twystd-qwerty
mqtt.topic.replies = /uiop
mqtt.topic.events = ./asdf
mqtt.topic.system = sys

# DEVICES
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.405419896.door.1 = Front Door
UT0311-L0x.405419896.door.2 = Side Door
UT0311-L0x.405419896.door.3 = Garage
UT0311-L0x.405419896.door.4 = Workshop
`)

func TestDefaultConfig(t *testing.T) {
	bind, broadcast, listen := DefaultIpAddresses()

	expected := Config{
		System: System{
			BindAddress:         &bind,
			BroadcastAddress:    &broadcast,
			ListenAddress:       &listen,
			HealthCheckInterval: 15 * time.Second,
			HealthCheckIdle:     60 * time.Second,
			HealthCheckIgnore:   5 * time.Minute,
			WatchdogInterval:    5 * time.Second,
		},

		MQTT: MQTT{
			Connection: Connection{
				ClientID: "twystd-uhppoted-mqttd",
			},
		},
	}

	config := NewConfig()

	if !reflect.DeepEqual(config.System, expected.System) {
		t.Errorf("Incorrect system default configuration:\nexpected:%+v,\ngot:     %+v", expected.System, config.System)
	}

	if config.MQTT.Connection.ClientID != expected.MQTT.Connection.ClientID {
		t.Errorf("Expected mqtt.connection.client.ID: '%v', got: '%v'", expected.MQTT.Connection.ClientID, config.MQTT.Connection.ClientID)
	}
}

func TestUnmarshal(t *testing.T) {
	expected := Config{
		System: System{
			BindAddress:         &net.UDPAddr{[]byte{192, 168, 1, 100}, 54321, ""},
			BroadcastAddress:    &net.UDPAddr{[]byte{192, 168, 1, 255}, 30000, ""},
			ListenAddress:       &net.UDPAddr{[]byte{192, 168, 1, 100}, 12345, ""},
			HealthCheckInterval: 31 * time.Second,
			HealthCheckIdle:     67 * time.Second,
			HealthCheckIgnore:   97 * time.Second,
			WatchdogInterval:    23 * time.Second,
		},

		MQTT: MQTT{
			Connection: Connection{
				Broker:            "tls://127.0.0.63:8887",
				ClientID:          "muppet",
				BrokerCertificate: "mqtt-broker.cert",
				ClientCertificate: "mqtt-client.cert",
				ClientKey:         "mqtt-client.key",
			},
		},
	}

	config := NewConfig()
	err := conf.Unmarshal(configuration, config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(config.System, expected.System) {
		t.Errorf("Incorrect system configuration:\nexpected:%+v,\ngot:     %+v", expected.System, config.System)
	}

	if !reflect.DeepEqual(config.Connection, expected.Connection) {
		t.Errorf("Incorrect 'mqtt.connection' configuration:\nexpected:%+v,\ngot:     %+v", expected.Connection, config.Connection)
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
		address := net.UDPAddr{
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
