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

# MQTT
mqtt.broker = 127.0.0.63:1887
mqtt.topic = twystd-qwerty

# DEVICES
UT0311-L0x.305419896.address = 192.168.1.100:60000
UT0311-L0x.305419896.door.1 = Front Door
UT0311-L0x.305419896.door.2 = Side Door
UT0311-L0x.305419896.door.3 = Garage
UT0311-L0x.305419896.door.4 = Workshop
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
		IP:   []byte{127, 0, 0, 63},
		Port: 1887,
		Zone: "",
	}

	if !reflect.DeepEqual(config.Broker, &address) {
		t.Errorf("Expected 'mqtt.broker' %s, got:%v", &address, config.Broker)
	}

	if config.Topic != "twystd-qwerty" {
		t.Errorf("Expected 'mqtt::topic' %v, got:%v", "twystd-qwerty", config.Topic)
	}

	if d, _ := config.Devices[305419896]; d == nil {
		t.Errorf("Expected 'device' for ID '%v', got:'%v'", 305419896, d)
	} else {
		address = net.UDPAddr{
			IP:   []byte{192, 168, 1, 100},
			Port: 60000,
			Zone: "",
		}

		if !reflect.DeepEqual(d.Address, &address) {
			t.Errorf("Expected 'device.address' %s for ID '%v', got:'%v'", &address, 305419896, d.Address)
		}

		if len(d.Door) != 4 {
			t.Errorf("Expected 4 entries for 'device.door' %s for ID '%v', got:%v", &address, 305419896, len(d.Door))
		} else {
			if d.Door[0] != "Front Door" {
				t.Errorf("Expected 'device.door[0]' %s for ID '%v', got:'%s'", "Front Door", 305419896, d.Door[0])
			}

			if d.Door[1] != "Side Door" {
				t.Errorf("Expected 'device.door[1]' %s for ID '%v', got:'%s'", "Side Door", 305419896, d.Door[1])
			}

			if d.Door[2] != "Garage" {
				t.Errorf("Expected 'device.door[2]' %s for ID '%v', got:'%s'", "Garage", 305419896, d.Door[2])
			}

			if d.Door[3] != "Workshop" {
				t.Errorf("Expected 'device.door[3]' %s for ID '%v', got:'%s'", "Workshop", 305419896, d.Door[3])
			}
		}
	}
}
