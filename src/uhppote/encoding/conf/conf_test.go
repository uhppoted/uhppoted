package conf

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type System struct {
	SysName string `conf:"sysname"`
	SysID   uint   `conf:"sysid"`
}

type Embedded struct {
	Name string `conf:"name"`
	ID   uint   `conf:"id"`
}

type testKV struct {
	key   string
	value interface{}
}

type testType struct {
	value string
}

func (f testType) MapKV(tag string, g func(string, interface{}) bool) bool {
	return g(tag, f.value)
}

func (f testType) MarshalConf(tag string) ([]byte, error) {
	var s strings.Builder

	fmt.Fprintf(&s, "%s = %s", tag, f.value)

	return []byte(s.String()), nil
}

func (f *testType) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	if v, ok := values[tag]; ok {
		return &testType{v}, nil
	}

	return f, nil
}

var configuration = []byte(
	`sysname = this
sysid = 42
udp.address = 192.168.1.100:54321
interface.value = qwerty
interface.pointer = uiop
sys.enabled = true
sys.byte = 127
sys.integer = -13579
sys.unsigned = 8081
sys.unsigned16 = 65535
sys.unsigned32 = 4294967295
sys.unsigned64 = 18446744073709551615
sys.string = asdfghjkl
sys.duration = 23s
embedded.name = zxcvb
embedded.id = 67890
`)

func TestMarshal(t *testing.T) {
	address := net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 54321,
		Zone: "",
	}

	config := struct {
		System     System
		Ignore     string        `conf:"-"`
		UdpAddress *net.UDPAddr  `conf:"udp.address"`
		Interface  testType      `conf:"interface.value"`
		InterfaceP *testType     `conf:"interface.pointer"`
		Enabled    bool          `conf:"sys.enabled"`
		Byte       byte          `conf:"sys.byte"`
		Integer    int           `conf:"sys.integer"`
		Unsigned   uint          `conf:"sys.unsigned"`
		Unsigned16 uint16        `conf:"sys.unsigned16"`
		Unsigned32 uint32        `conf:"sys.unsigned32"`
		Unsigned64 uint64        `conf:"sys.unsigned64"`
		String     string        `conf:"sys.string"`
		Duration   time.Duration `conf:"sys.duration"`
		Embedded   `conf:"embedded"`
	}{
		System: System{
			SysName: "this",
			SysID:   42,
		},
		Ignore:     "ignore",
		UdpAddress: &address,
		Interface:  testType{"qwerty"},
		InterfaceP: &testType{"uiop"},
		Enabled:    true,
		Byte:       127,
		Integer:    -13579,
		Unsigned:   8081,
		Unsigned16: 65535,
		Unsigned32: 4294967295,
		Unsigned64: 18446744073709551615,
		String:     "asdfghjkl",
		Duration:   23 * time.Second,
		Embedded: Embedded{
			Name: "zxcvb",
			ID:   67890,
		},
	}

	bytes, err := Marshal(config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(bytes, configuration) {
		l, ls, p, q := diff(string(configuration), string(bytes))
		t.Errorf("conf not marshaled correctly:\n%s\n>> line %d:\n>> %s\n--------\n   %s\n   %s\n--------\n", string(bytes), l, ls, p, q)
	}
}

func TestUnmarshal(t *testing.T) {
	config := struct {
		System     System        `conf:""`
		Ignore     string        `conf:"-"`
		UdpAddress *net.UDPAddr  `conf:"udp.address"`
		Interface  testType      `conf:"interface.value"`
		InterfaceP *testType     `conf:"interface.pointer"`
		Enabled    bool          `conf:"sys.enabled"`
		Byte       byte          `conf:"sys.byte"`
		Integer    int           `conf:"sys.integer"`
		Unsigned   uint          `conf:"sys.unsigned"`
		Unsigned16 uint16        `conf:"sys.unsigned16"`
		Unsigned32 uint32        `conf:"sys.unsigned32"`
		Unsigned64 uint64        `conf:"sys.unsigned64"`
		String     string        `conf:"sys.string"`
		Duration   time.Duration `conf:"sys.duration"`
		Embedded   `conf:"embedded"`
	}{
		Ignore: "ignore",
	}

	err := Unmarshal(configuration, &config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	address := net.UDPAddr{
		IP:   []byte{192, 168, 1, 100},
		Port: 54321,
		Zone: "",
	}

	if config.System.SysName != "this" {
		t.Errorf("Expected 'sysname' %s, got: %s", "this", config.System.SysName)
	}

	if config.System.SysID != 42 {
		t.Errorf("Expected 'sysid' %d, got: %d", 42, config.System.SysID)
	}

	if config.Ignore != "ignore" {
		t.Errorf("Expected 'ignore' %s, got: %s", "ignore", config.Ignore)
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

	if config.Byte != byte(127) {
		t.Errorf("Expected 'byte' value '%v', got: '%v'", 127, config.Byte)
	}

	if config.Integer != -13579 {
		t.Errorf("Expected 'integer' value '%v', got: '%v'", -13579, config.Integer)
	}

	if config.Unsigned != 8081 {
		t.Errorf("Expected 'unsigned' value '%v', got: '%v'", 8081, config.Unsigned)
	}

	if config.Unsigned16 != uint16(65535) {
		t.Errorf("Expected 'unsigned16' value '%v', got: '%v'", uint16(65535), config.Unsigned16)
	}

	if config.Unsigned32 != uint32(4294967295) {
		t.Errorf("Expected 'unsigned32' value '%v', got: '%v'", uint32(4294967295), config.Unsigned32)
	}

	if config.Unsigned64 != uint64(18446744073709551615) {
		t.Errorf("Expected 'unsigned64' value '%v', got: '%v'", uint64(18446744073709551615), config.Unsigned64)
	}

	if config.String != "asdfghjkl" {
		t.Errorf("Expected 'string' value '%v', got: '%v'", "asdfghjkl", config.String)
	}

	if config.Duration != 23*time.Second {
		t.Errorf("Expected 'duration' value '%v', got: '%v'", 23*time.Second, config.Duration)
	}

	if config.Name != "zxcvb" {
		t.Errorf("Expected 'embedded.name' value '%v', got: '%v'", "zxcvb", config.Name)
	}

	if config.ID != uint(67890) {
		t.Errorf("Expected 'embedded.id' value '%v', got: '%v'", uint(67890), config.ID)
	}
}

func TestRange(t *testing.T) {
	address, _ := net.ResolveUDPAddr("udp", "192.168.1.100:54321")

	config := struct {
		Ignore     string        `conf:"-"`
		UdpAddress *net.UDPAddr  `conf:"udp.address"`
		Interface  testType      `conf:"interface.value"`
		InterfaceP *testType     `conf:"interface.pointer"`
		Enabled    bool          `conf:"sys.enabled"`
		Byte       byte          `conf:"sys.byte"`
		Integer    int           `conf:"sys.integer"`
		Unsigned   uint          `conf:"sys.unsigned"`
		Unsigned16 uint16        `conf:"sys.unsigned16"`
		Unsigned32 uint32        `conf:"sys.unsigned32"`
		Unsigned64 uint64        `conf:"sys.unsigned64"`
		String     string        `conf:"sys.string"`
		Duration   time.Duration `conf:"sys.duration"`
		Embedded   `conf:"embedded"`
	}{
		Ignore:     "ignore",
		UdpAddress: address,
		Interface:  testType{"qwerty"},
		InterfaceP: &testType{"uiop"},
		Enabled:    true,
		Byte:       127,
		Integer:    -13579,
		Unsigned:   8081,
		Unsigned16: 65535,
		Unsigned32: 4294967295,
		Unsigned64: 18446744073709551615,
		String:     "asdfghjkl",
		Duration:   23 * time.Second,
		Embedded: Embedded{
			Name: "zxcvb",
			ID:   67890,
		},
	}

	expected := []testKV{
		testKV{"udp.address", address},
		testKV{"interface.value", "qwerty"},
		testKV{"interface.pointer", "uiop"},
		testKV{"sys.enabled", true},
		testKV{"sys.byte", byte(127)},
		testKV{"sys.integer", -13579},
		testKV{"sys.unsigned", uint(8081)},
		testKV{"sys.unsigned16", uint16(65535)},
		testKV{"sys.unsigned32", uint32(4294967295)},
		testKV{"sys.unsigned64", uint64(18446744073709551615)},
		testKV{"sys.string", "asdfghjkl"},
		testKV{"sys.duration", 23 * time.Second},
		testKV{"embedded.name", "zxcvb"},
		testKV{"embedded.id", uint(67890)},
	}

	list := []testKV{}
	Range(config, func(k string, v interface{}) bool {
		list = append(list, testKV{k, v})
		return true
	})

	if !reflect.DeepEqual(list, expected) {
		var err strings.Builder
		fmt.Fprintf(&err, "Range did not fill list correctly:\n\n--------\n   %v\n   %v\n\n", expected, list)

		for ix := 0; ix < len(expected) && ix < len(list); ix++ {
			if !reflect.DeepEqual(list[ix], expected[ix]) {
				fmt.Fprintf(&err, "   %#v\n   %#v\n", expected[ix], list[ix])
				break
			}
		}
		fmt.Fprintf(&err, "--------\n")

		t.Errorf(err.String())
	}
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
