package config

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Device struct {
	Address *net.UDPAddr
	Door1   string
	Door2   string
	Door3   string
	Door4   string
}

type Config struct {
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
	Devices          map[uint32]*Device
}

var parsers = []struct {
	re *regexp.Regexp
	f  func(string, *Config) *Config
}{
	{regexp.MustCompile("^bind\\.address\\s*=.*"), bind},
	{regexp.MustCompile("^broadcast\\.address\\s*=.*"), broadcast},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.address\\s*=.*"), address},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.door\\.[1-4]\\s*=.*"), door},
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	c := &Config{
		BindAddress:      nil,
		BroadcastAddress: nil,
		Devices:          make(map[uint32]*Device),
	}

	for s.Scan() {
		l := s.Text()
		for _, p := range parsers {
			if p.re.MatchString(l) {
				c = p.f(l, c)
			}
		}
	}

	return c, nil
}

func bind(l string, c *Config) *Config {
	re := regexp.MustCompile("^bind\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[1]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid UDP bind address '%s': %v\n", l, err)
		} else {
			c.BindAddress = address
		}
	}

	return c
}

func broadcast(l string, c *Config) *Config {
	re := regexp.MustCompile("^broadcast\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[1]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid UDP broadcast address '%s': %v\n", l, err)
		} else {
			c.BroadcastAddress = address
		}
	}

	return c
}

func address(l string, c *Config) *Config {
	re := regexp.MustCompile("^UT0311-L0x\\.([0-9]+)\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		serialNo, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid serial number '%s': %v\n", l, err)
			return c
		}

		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[2]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid device UDP address '%s': %v\n", l, err)
			return c
		}

		k := uint32(serialNo)
		d := c.Devices[k]
		if d == nil {
			d = &Device{}
		}

		d.Address = address
		c.Devices[k] = d
	}

	return c
}

func door(l string, c *Config) *Config {
	re := regexp.MustCompile("^UT0311-L0x\\.([0-9]+)\\.door\\.([1-4])\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		serialNo, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid serial number '%s': %v\n", l, err)
			return c
		}

		door, err := strconv.ParseUint(match[2], 10, 8)
		if err != nil || door < 1 || door > 4 {
			fmt.Printf("WARN: configuration error - invalid device door number '%s': %v\n", l, err)
			return c
		}

		k := uint32(serialNo)
		d := c.Devices[k]
		if d == nil {
			d = &Device{}
		}

		name := strings.TrimSpace(match[3])
		switch door {
		case 1:
			d.Door1 = name
		case 2:
			d.Door2 = name
		case 3:
			d.Door3 = name
		case 4:
			d.Door4 = name
		}

		c.Devices[k] = d
	}

	return c
}
