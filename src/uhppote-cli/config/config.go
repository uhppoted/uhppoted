package config

import (
	"bufio"
	"encoding/binary"
	"errors"
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
	f  func(string, *Config) *Config
}{
	{regexp.MustCompile("^bind\\.address\\s*=.*"), bind},
	{regexp.MustCompile("^broadcast\\.address\\s*=.*"), broadcast},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.address\\s*=.*"), address},
	{regexp.MustCompile("^UT0311-L0x\\.[0-9]+\\.door\\.[1-4]\\s*=.*"), door},
}

func LoadConfig(path string) (*Config, error) {
	bind, broadcast, err := DefaultIpAddresses()
	if err != nil {
		return nil, err
	}

	if bind == nil || broadcast == nil {
		return nil, errors.New("Unable to determine default bind and broadcast IP addresses")
	}

	c := &Config{
		File:             path,
		BindAddress:      *bind,
		BroadcastAddress: *broadcast,
		Devices:          make(map[uint32]*Device),
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
				c = p.f(l, c)
			}
		}
	}

	return c, nil
}

func (c *Config) Verify() error {
	doors := make(map[string]bool)
	for _, device := range c.Devices {
		for _, door := range device.Door {
			d := strings.ReplaceAll(strings.ToLower(door), " ", "")

			if doors[d] {
				return errors.New(fmt.Sprintf("Door '%s' is defined more than once in configuration file '%s'", door, c.File))
			}

			doors[d] = true
		}
	}

	return nil
}

func bind(l string, c *Config) *Config {
	re := regexp.MustCompile("^bind\\.address\\s*=(.*)")
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		address, err := net.ResolveUDPAddr("udp", strings.TrimSpace(match[1]))
		if err != nil {
			fmt.Printf("WARN: configuration error - invalid UDP bind address '%s': %v\n", l, err)
		} else if address == nil {
			fmt.Println("WARN: configuration error - invalid UDP bind address")
		} else {
			c.BindAddress = *address
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
		} else if address == nil {
			fmt.Println("WARN: configuration error - invalid UDP broadcast address")
		} else {
			c.BroadcastAddress = *address
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
			d = &Device{Door: make([]string, 4)}
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
			d = &Device{Door: make([]string, 4)}
		}

		d.Door[door-1] = strings.TrimSpace(match[3])
		c.Devices[k] = d
	}

	return c
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (*net.UDPAddr, *net.UDPAddr, error) {
	bind := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 0,
		Zone: "",
	}

	broadcast := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60000,
		Zone: "",
	}

	copy(bind.IP, net.IPv4zero)
	copy(broadcast.IP, net.IPv4bcast)

	if ifaces, err := net.Interfaces(); err == nil {
	loop:
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err == nil {
				for _, a := range addrs {
					switch v := a.(type) {
					case *net.IPNet:
						if v.IP.To4() != nil && i.Flags&net.FlagLoopback == 0 {
							copy(bind.IP, v.IP.To4())
							if i.Flags&net.FlagBroadcast != 0 {
								addr := v.IP.To4()
								mask := v.Mask
								binary.BigEndian.PutUint32(broadcast.IP, binary.BigEndian.Uint32(addr)|^binary.BigEndian.Uint32(mask))
							}
							break loop
						}
					}
				}
			}
		}
	}

	return &bind, &broadcast, nil
}
