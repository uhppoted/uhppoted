package config

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Device struct {
	Address *net.UDPAddr
	Door    []string
}

type REST struct {
	HttpEnabled        bool
	HttpPort           uint16
	HttpsEnabled       bool
	HttpsPort          uint16
	TLSKeyFile         string
	TLSCertificateFile string
	CACertificateFile  string
	CORSEnabled        bool
}

type OpenApi struct {
	Enabled   bool
	Directory string
}

type Config struct {
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
	Devices          map[uint32]*Device
	REST
	OpenApi
}

var parsers = []struct {
	re *regexp.Regexp
	f  func(string, *Config)
}{
	{regexp.MustCompile(`^bind\.address\s*=.*`), bind},
	{regexp.MustCompile(`^broadcast\.address\s*=.*`), broadcast},
	{regexp.MustCompile(`^UT0311-L0x\.[0-9]+\.address\s*=.*`), address},
	{regexp.MustCompile(`^UT0311-L0x\.[0-9]+\.door\.[1-4]\s*=.*`), door},
	{regexp.MustCompile(`^rest\.http\.enabled\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.http\.port\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.https\.enabled\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.https\.port\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.tls\.key\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.tls\.certificate\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.tls\.ca\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.CORS\.enabled\s*=.*`), rest},
	{regexp.MustCompile(`^rest\.openapi\.enabled\s*=.*`), openapi},
	{regexp.MustCompile(`^rest\.openapi\.directory\s*=.*`), openapi},
}

func NewConfig() *Config {
	bind, broadcast := DefaultIpAddresses()

	c := Config{
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
		Devices:          make(map[uint32]*Device),

		REST: REST{
			HttpEnabled:        false,
			HttpPort:           8080,
			HttpsEnabled:       true,
			HttpsPort:          8443,
			TLSKeyFile:         "uhppoted.key",
			TLSCertificateFile: "uhppoted.cert",
			CACertificateFile:  "ca.cert",
			CORSEnabled:        false,
		},

		OpenApi: OpenApi{
			Enabled:   false,
			Directory: "./openapi",
		},
	}

	return &c
}

func (c *Config) Load(path string) error {
	if path == "" {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	c.OpenApi.Directory = filepath.Join(filepath.Dir(path), "rest", "openapi")

	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		for _, p := range parsers {
			if p.re.MatchString(l) {
				p.f(l, c)
			}
		}
	}

	return nil
}

func bind(l string, c *Config) {
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
}

func broadcast(l string, c *Config) {
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
	re := regexp.MustCompile(`^UT0311-L0x\.([0-9]+)\.door\.([1-4])\s*=(.*)`)
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

func rest(l string, c *Config) {
	re := regexp.MustCompile(`^rest\.(\w+)\.(\w+)\s*=(.*)`)
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		switch match[1] {
		case "http":
			switch match[2] {
			case "enabled":
				if enabled, err := strconv.ParseBool(strings.TrimSpace(match[3])); err == nil {
					c.REST.HttpEnabled = enabled
				}
			case "port":
				if port, err := strconv.ParseUint(strings.TrimSpace(match[3]), 10, 16); err == nil {
					c.REST.HttpPort = uint16(port)
				}
			}

		case "https":
			switch match[2] {
			case "enabled":
				if enabled, err := strconv.ParseBool(strings.TrimSpace(match[3])); err == nil {
					c.REST.HttpsEnabled = enabled
				}
			case "port":
				if port, err := strconv.ParseUint(strings.TrimSpace(match[3]), 10, 16); err == nil {
					c.REST.HttpsPort = uint16(port)
				}
			}

		case "tls":
			switch match[2] {
			case "key":
				c.REST.TLSKeyFile = strings.TrimSpace(match[3])
			case "certificate":
				c.REST.TLSCertificateFile = strings.TrimSpace(match[3])
			case "ca":
				c.REST.CACertificateFile = strings.TrimSpace(match[3])
			}

		case "CORS":
			switch match[2] {
			case "enabled":
				if enabled, err := strconv.ParseBool(strings.TrimSpace(match[3])); err == nil {
					c.REST.CORSEnabled = enabled
				}
			}
		}
	}
}

func openapi(l string, c *Config) {
	re := regexp.MustCompile(`^rest\.openapi\.(\w+)\s*=(.*)`)
	match := re.FindStringSubmatch(l)

	if len(match) > 0 {
		switch match[1] {
		case "enabled":
			if enabled, err := strconv.ParseBool(strings.TrimSpace(match[2])); err == nil {
				c.OpenApi.Enabled = enabled
			}

		case "directory":
			c.OpenApi.Directory = strings.TrimSpace(match[2])
		}
	}
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (net.UDPAddr, net.UDPAddr) {
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

	return bind, broadcast
}
