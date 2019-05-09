package config

import (
	"net"
)

type Config struct {
	Devices map[uint32]*net.UDPAddr
}

func NewConfig() (*Config, error) {
	devices := make(map[uint32]*net.UDPAddr)

	devices[423187757], _ = net.ResolveUDPAddr("udp", "192.168.0.14:60000")

	return &Config{
		Devices: devices,
	}, nil
}
