package config

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

type Device struct {
	Address *net.UDPAddr
}

type Config struct {
	Devices map[uint32]*Device
}

var re = regexp.MustCompile("^UT0311-L0x\\.([0-9]+)\\.(.*)")

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	m := make(map[uint32]*Device)
	s := bufio.NewScanner(f)

	fmt.Println()
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())

		if len(match) > 0 {
			serialNumber, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				fmt.Printf("WARN: configuration error - invalid serial number '%s': %v\n", s.Text(), err)
				continue
			} else {
				d := m[uint32(serialNumber)]
				if d == nil {
					d = &Device{
						Address: nil,
					}
				}
				m[uint32(serialNumber)] = parse(d, match[2])
			}
		}
	}

	return &Config{
		Devices: m,
	}, nil
}

func parse(d *Device, attr string) *Device {
	// Extract device UDP address

	re := regexp.MustCompile("\\s*address\\s*=\\s*([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]+)")
	match := re.FindStringSubmatch(attr)

	if len(match) > 1 {
		address, err := net.ResolveUDPAddr("udp", match[1])
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid address '%s': %v\n", attr, err)
		} else {
			d.Address = address
		}

	}

	return d
}
