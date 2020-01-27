package conf

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type testType struct {
	value string
}

type Embedded struct {
	Name string `conf:"name"`
	ID   uint   `conf:"id"`
}

var configuration = []byte(
	`udp.address = 192.168.1.100:54321
interface.value = qwerty
interface.pointer = uiop
sys.enabled = true
sys.integer = -13579
sys.unsigned = 8081
sys.string = asdfghjkl
embedded.name = zxcvb
embedded.id = 67890
`)

func TestMarshal(t *testing.T) {
	expected := `udp.address = 192.168.1.100:54321
interface.value = qwerty
interface.pointer = uiop
sys.enabled = true
sys.integer = -13579
sys.unsigned = 8081
sys.string = asdfghjkl
embedded.name = zxcvb
embedded.id = 67890
`

	address := net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 54321,
		Zone: "",
	}

	config := struct {
		UdpAddress *net.UDPAddr `conf:"udp.address"`
		Interface  testType     `conf:"interface.value"`
		InterfaceP *testType    `conf:"interface.pointer"`
		Enabled    bool         `conf:"sys.enabled"`
		Integer    int          `conf:"sys.integer"`
		Unsigned   uint         `conf:"sys.unsigned"`
		String     string       `conf:"sys.string"`
		Embedded   `conf:"embedded"`
	}{
		UdpAddress: &address,
		Interface:  testType{"qwerty"},
		InterfaceP: &testType{"uiop"},
		Enabled:    true,
		Integer:    -13579,
		Unsigned:   8081,
		String:     "asdfghjkl",
		Embedded: Embedded{
			Name: "zxcvb",
			ID:   67890,
		},
	}

	bytes, err := Marshal(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(bytes) != expected {
		l, ls, p, q := diff(expected, string(bytes))
		t.Errorf("conf not marshaled correctly:\n%s\n>> line %d:\n>> %s\n--------\n   %s\n   %s\n--------\n", string(bytes), l, ls, p, q)
	}
}

func TestUnmarshal(t *testing.T) {
	config := struct {
		UdpAddress *net.UDPAddr `conf:"udp.address"`
		Interface  testType     `conf:"interface.value"`
		InterfaceP *testType    `conf:"interface.pointer"`
		Enabled    bool         `conf:"sys.enabled"`
		Integer    int          `conf:"sys.integer"`
		Unsigned   uint         `conf:"sys.unsigned"`
		String     string       `conf:"sys.string"`
		Embedded   `conf:"embedded"`
	}{}

	err := Unmarshal(configuration, &config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	address := net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 54321,
		Zone: "",
	}

	if !reflect.DeepEqual(config.UdpAddress, &address) {
		t.Errorf("Expected 'udp.address' %s, got: %s", &address, config.UdpAddress)
	}

	if config.Interface.value != "qwerty" {
		t.Errorf("Expected 'interface' value '%s', got: '%v'", "qwerty", config.Interface)
	}

	if config.InterfaceP == nil || config.InterfaceP.value != "uiop" {
		t.Errorf("Expected 'interface pointer' value '%s', got: '%v'", "uiop", config.InterfaceP)
	}

	if !config.Enabled {
		t.Errorf("Expected 'boolean' value '%v', got: '%v'", true, config.Enabled)
	}

	if config.Integer != -13579 {
		t.Errorf("Expected 'integer' value '%v', got: '%v'", -13579, config.Integer)
	}

	if config.Unsigned != 8081 {
		t.Errorf("Expected 'unsigned' value '%v', got: '%v'", 8081, config.Unsigned)
	}

	if config.String != "asdfghjkl" {
		t.Errorf("Expected 'string' value '%v', got: '%v'", "asdfghjkl", config.String)
	}

	if config.Name != "zxcvb" {
		t.Errorf("Expected 'embedded.name' value '%v', got: '%v'", "zxcvb", config.Name)
	}

	if config.ID != 67890 {
		t.Errorf("Expected 'embedded.id' value '%v', got: '%v'", 67890, config.ID)
	}
}

func (f testType) MarshalConf() ([]byte, error) {
	return []byte(f.value), nil
}

func (f *testType) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	if v, ok := values[tag]; ok {
		return &testType{v}, nil
	}

	return f, nil
}

// Unmarshal example for map[id]device using Unmarshaler interface

type deviceMap map[uint32]*device

type device struct {
	name    string
	address string
}

func (d device) String() string {
	return fmt.Sprintf("%-7s %s", d.name, d.address)
}

func ExampleUnmarshal() {
	configuration := `# DEVICES
UT0311-L0x.405419896.name = BOARD1
UT0311-L0x.405419896.address = 192.168.1.100:60000
UT0311-L0x.54321.name = BOARD2
UT0311-L0x.54321.address = 192.168.1.101:60000
`

	config := struct {
		Devices deviceMap `conf:"/UT0311-L0x\\.([0-9]+)\\.(\\w+)/"`
	}{}

	err := Unmarshal([]byte(configuration), &config)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return
	}

	for id, d := range config.Devices {
		fmt.Printf("DEVICE: %-10d %s\n", id, d)
	}

	// Unordered output:
	// DEVICE: 405419896  BOARD1  192.168.1.100:60000
	// DEVICE: 54321      BOARD2  192.168.1.101:60000
}

func (f *deviceMap) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(`^/(.*?)/$`)
	match := re.FindStringSubmatch(tag)
	if len(match) < 2 {
		return f, fmt.Errorf("Invalid 'conf' regular expression tag: %s", tag)
	}

	re, err := regexp.Compile(match[1])
	if err != nil {
		return f, err
	}

	var m deviceMap

	if f != nil {
		m = *f
	}

	if m == nil {
		m = make(deviceMap, 0)
	}

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) == 3 {
			id, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				return f, fmt.Errorf("Invalid 'deviceMap' key %s: %v", key, err)
			}

			d, ok := m[uint32(id)]
			if !ok || d == nil {
				d = &device{}
				m[uint32(id)] = d
			}

			switch match[2] {
			case "name":
				d.name = value
			case "address":
				d.address = value
			}
		}
	}

	return &m, nil
}

func diff(p, q string) (int, string, string, string) {
	line := 0
	s := strings.Split(p, "\n")
	t := strings.Split(q, "\n")

	for ix := 0; ix < len(s) && ix < len(t); ix++ {
		line++
		println(">>", s[ix])
		println(">>", t[ix])
		println()
		if s[ix] != t[ix] {
			u := []rune(s[ix])
			v := []rune(t[ix])
			for jx := 0; jx < len(u) && jx < len(v); jx++ {
				if u[jx] != v[jx] {
					break
				}
				u[jx] = '.'
				v[jx] = '.'
			}

			return line, s[ix], string(u), string(v)
		}
	}

	return line, "?", "?", "?"
}
