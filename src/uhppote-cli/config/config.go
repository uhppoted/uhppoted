package config

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Config struct {
	Devices map[uint32]*net.UDPAddr
}

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	fmt.Println()
	for s.Scan() {
		fmt.Printf(" >>> DEBUG: %s\n", s.Text())
	}
	fmt.Println()

	devices := make(map[uint32]*net.UDPAddr)

	devices[423187757], _ = net.ResolveUDPAddr("udp", "192.168.0.14:60000")

	return &Config{
		Devices: devices,
	}, nil
}
