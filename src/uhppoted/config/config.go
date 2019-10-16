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
	Door    []string
}

type Config struct {
	File             string
	BindAddress      net.UDPAddr
	BroadcastAddress net.UDPAddr
	Devices          map[uint32]*Device
}

var parsers = []struct {
	re *regexp.Regexp
	f  func(string, *Config)
}{
	{regexp.MustCompile("^bind\\.address\\s*=.*"), bind},
	{regexp.MustCompile("^broadcast\\.address\\s*=.*"), broadcast},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.address\\s*=.*"), address},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.door\\.[1-4]\\s*=.*"), door},
}

func LoadConfig(path string) (Config, error) {
	c := Config{
		File: path,
		BindAddress: bindAddr(),
		BroadcastAddress: broadcastAddr(),
		Devices: make(map[uint32]*Device),
	}

	if path == "" {
		return c, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return c, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		for _, p := range parsers {
			if p.re.MatchString(l) {
				p.f(l, &c)
			}
		}
	}

	return c, nil
}

func bind(l string, c *Config) {
	re := regexp.MustCompile("^bind\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[1]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid UDP bind address '%s': %v\n", l, err)
		} else {
			c.BindAddress = *address
		}
	}
}

func broadcast(l string, c *Config){
	re := regexp.MustCompile("^broadcast\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[1]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid UDP broadcast address '%s': %v\n", l, err)
		} else {
			c.BroadcastAddress = *address
		}
	}
}

func address(l string, c *Config) {
	re := regexp.MustCompile("^UT0311-L0x\\.([0-9]+)\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		serialNo, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid serial number '%s': %v\n", l, err)
			return
		}

		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[2]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid device UDP address '%s': %v\n", l, err)
			return
		}

		k := uint32(serialNo)
		d := c.Devices[k]
		if d == nil {
			d = &Device{Door: make([]string, 4)}
		}

		d.Address = address
		c.Devices[k] = d
	}
}

func door(l string, c *Config) {
	re := regexp.MustCompile("^UT0311-L0x\\.([0-9]+)\\.door\\.([1-4])\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		serialNo, err := strconv.ParseUint(match[1], 10, 32)
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid serial number '%s': %v\n", l, err)
			return
		}

		door, err := strconv.ParseUint(match[2], 10, 8)
		if err != nil || door < 1 || door > 4 {
			fmt.Printf("WARN: configuration error - invalid device door number '%s': %v\n", l, err)
			return
		}

		k := uint32(serialNo)
		d := c.Devices[k]
		if d == nil {
			d = &Device{Door: make([]string, 4)}
		}

		d.Door[door-1] = strings.TrimSpace(match[3])
		c.Devices[k] = d
	}
}

func bindAddr() net.UDPAddr {
	return net.UDPAddr{
		IP: net.IPv4zero,
		Port: 0,
		Zone: "",
		}
}

func broadcastAddr() net.UDPAddr {
	return net.UDPAddr{
		IP: net.IPv4bcast,
		Port: 60000,
		Zone: "",
		}
}
