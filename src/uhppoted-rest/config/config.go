package config

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"uhppote/encoding/conf"
)

type DeviceMap map[uint32]*Device

type Device struct {
	Address  *net.UDPAddr
	Rollover uint32
	Door     []string
}

type REST struct {
	HttpEnabled        bool   `conf:"http.enabled"`
	HttpPort           uint16 `conf:"http.port"`
	HttpsEnabled       bool   `conf:"https.enabled"`
	HttpsPort          uint16 `conf:"https.port"`
	TLSKeyFile         string `conf:"tls.key"`
	TLSCertificateFile string `conf:"tls.certificate"`
	CACertificateFile  string `conf:"tls.ca"`
	CORSEnabled        bool   `conf:"CORS.enabled"`
}

type OpenApi struct {
	Enabled   bool   `conf:"enabled"`
	Directory string `conf:"directory"`
}

type Config struct {
	BindAddress      *net.UDPAddr `conf:"bind.address"`
	BroadcastAddress *net.UDPAddr `conf:"broadcast.address"`
	Devices          DeviceMap    `conf:"/^UT0311-L0x\\.([0-9]+)\\.(.*)/"`
	REST             `conf:"rest"`
	OpenApi          `conf:"openapi"`
}

const ROLLOVER = 100000

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

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	return conf.Unmarshal(bytes, c)
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

func (f *DeviceMap) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			id, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return f, fmt.Errorf("Invalid 'testMap' key %s: %v", key, err)
			}

			d, ok := (*f)[uint32(id)]
			if !ok || d == nil {
				d = &Device{
					Door:     make([]string, 4),
					Rollover: ROLLOVER,
				}

				(*f)[uint32(id)] = d
			}

			switch match[2] {
			case "address":
				address, err := net.ResolveUDPAddr("udp", value)
				if err != nil {
					return f, fmt.Errorf("Device %v, invalid address '%s': %v", id, value, err)
				} else {
					d.Address = &net.UDPAddr{
						IP:   make(net.IP, net.IPv4len),
						Port: address.Port,
						Zone: "",
					}

					copy(d.Address.IP, address.IP.To4())
				}

			case "rollover":
				rollover, err := strconv.ParseUint(strings.TrimSpace(value), 10, 32)
				if err != nil {
					return f, fmt.Errorf("Device %v, invalid rollover '%s': %v", id, value, err)
				} else {
					d.Rollover = uint32(rollover)
				}

			case "door.1":
				d.Door[0] = value

			case "door.2":
				d.Door[1] = value

			case "door.3":
				d.Door[2] = value

			case "door.4":
				d.Door[3] = value
			}
		}
	}

	return f, nil
}
