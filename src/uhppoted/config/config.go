package config

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	"uhppote/encoding/conf"
)

type DeviceMap map[uint32]*Device

type Device struct {
	Address *net.UDPAddr
	Door    []string
}

// # OPEN API
// # openapi.enabled = false
// # openapi.directory = {{.WorkDir}}\rest\openapi

type kv struct {
	Key   string
	Value interface{}
}

const pretty = `# SYSTEM{{range .system}}
{{.Key}} = {{.Value}}{{end}}

# REST{{range .rest}}
{{.Key}} = {{.Value}}{{end}}

# MQTT{{range .mqtt}}
{{.Key}} = {{.Value}}{{end}}

# OPEN API{{range .openapi}}
# {{.Key}} = {{.Value}}{{end}}

# DEVICES{{range $id,$device := .devices}}
UT0311-L0x.{{$id}}.address = {{$device.Address}}
UT0311-L0x.{{$id}}.door.1 = {{index $device.Door 0}}
UT0311-L0x.{{$id}}.door.2 = {{index $device.Door 1}}
UT0311-L0x.{{$id}}.door.3 = {{index $device.Door 2}}
UT0311-L0x.{{$id}}.door.4 = {{index $device.Door 3}}
{{else}}
# Example configuration for UTO311-L04 with serial number 405419896
# UT0311-L0x.405419896.address = 192.168.1.100:60000
# UT0311-L0x.405419896.door.1 = Front Door
# UT0311-L0x.405419896.door.2 = Side Door
# UT0311-L0x.405419896.door.3 = Garage
# UT0311-L0x.405419896.door.4 = Workshop
{{end}}`

type Config struct {
	BindAddress         *net.UDPAddr  `conf:"bind.address"`
	BroadcastAddress    *net.UDPAddr  `conf:"broadcast.address"`
	ListenAddress       *net.UDPAddr  `conf:"listen.address"`
	HealthCheckInterval time.Duration `conf:"monitoring.healthcheck.interval"`
	WatchdogInterval    time.Duration `conf:"monitoring.watchdog.interval"`
	Devices             DeviceMap     `conf:"/^UT0311-L0x\\.([0-9]+)\\.(.*)/"`
	REST                `conf:"rest"`
	MQTT                `conf:"mqtt"`
	OpenAPI             `conf:"openapi"`
}

func NewConfig() *Config {
	bind, broadcast, listen := DefaultIpAddresses()

	c := Config{
		BindAddress:         &bind,
		BroadcastAddress:    &broadcast,
		ListenAddress:       &listen,
		HealthCheckInterval: 15 * time.Second,
		WatchdogInterval:    5 * time.Second,
		REST:                *NewREST(),
		MQTT:                *NewMQTT(),
		OpenAPI:             *NewOpenAPI(),
		Devices:             make(DeviceMap, 0),
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

	return c.Read(f)
}

func (c *Config) Read(r io.Reader) error {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return conf.Unmarshal(bytes, c)
}

func (c *Config) Write(w io.Writer) error {
	system := []kv{}
	v := reflect.ValueOf(c)
	s := v.Elem()
	N := s.NumField()
	for i := 0; i < N; i++ {
		f := s.Field(i)
		t := s.Type().Field(i)
		tag := t.Tag.Get("conf")
		if f.Kind() != reflect.Struct && f.Kind() != reflect.Map {
			system = append(system, kv{tag, f})
		}
	}

	config := map[string]interface{}{
		"system":  system,
		"rest":    listify("rest", reflect.ValueOf(c.REST)),
		"mqtt":    listify("mqtt", reflect.ValueOf(c.MQTT)),
		"openapi": listify("openapi", reflect.ValueOf(c.OpenAPI)),
		"devices": c.Devices,
	}

	return template.Must(template.New("uhppoted.conf").Parse(pretty)).Execute(w, config)
}

func listify(parent string, s reflect.Value) []kv {
	list := []kv{}
	N := s.NumField()
	for i := 0; i < N; i++ {
		f := s.Field(i)
		t := s.Type().Field(i)
		tag := t.Tag.Get("conf")

		if f.Kind() == reflect.Struct {
			list = append(list, listify(parent+"."+tag, f)...)
		} else {
			list = append(list, kv{parent + "." + tag, fmt.Sprintf("%v", f)})
		}
	}

	return list
}

// Ref. https://stackoverflow.com/questions/23529663/how-to-get-all-addresses-and-masks-from-local-interfaces-in-go
func DefaultIpAddresses() (net.UDPAddr, net.UDPAddr, net.UDPAddr) {
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

	listen := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: 60001,
		Zone: "",
	}

	copy(bind.IP, net.IPv4zero)
	copy(broadcast.IP, net.IPv4bcast)
	copy(listen.IP, net.IPv4zero)

	if ifaces, err := net.Interfaces(); err == nil {
	loop:
		for _, i := range ifaces {
			if addrs, err := i.Addrs(); err == nil {
				for _, a := range addrs {
					switch v := a.(type) {
					case *net.IPNet:
						if v.IP.To4() != nil && i.Flags&net.FlagLoopback == 0 {
							copy(bind.IP, v.IP.To4())
							copy(listen.IP, v.IP.To4())
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

	return bind, broadcast, listen
}

func (f DeviceMap) MarshalConf(tag string) ([]byte, error) {
	var s strings.Builder

	if len(f) > 0 {
		fmt.Fprintf(&s, "# DEVICES\n")
		for id, device := range f {
			fmt.Fprintf(&s, "UTO311-L0x.%d.address = %s\n", id, device.Address)
			for d, door := range device.Door {
				fmt.Fprintf(&s, "UTO311-L0x.%d.door.%d = %s\n", id, d+1, door)
			}
			fmt.Fprintf(&s, "\n")
		}
	}

	return []byte(s.String()), nil
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
					Door: make([]string, 4),
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

func resolve(v string) (*net.UDPAddr, error) {
	address, err := net.ResolveUDPAddr("udp", v)
	if err != nil {
		return nil, err
	}

	addr := net.UDPAddr{
		IP:   make(net.IP, net.IPv4len),
		Port: address.Port,
		Zone: "",
	}

	copy(addr.IP, address.IP.To4())

	return &addr, nil
}
