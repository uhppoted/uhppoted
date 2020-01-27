package conf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Marshaler interface {
	MarshalConf() ([]byte, error)
}

type Unmarshaler interface {
	UnmarshalConf(tag string, values map[string]string) (interface{}, error)
}

var (
	tBool     = reflect.TypeOf(bool(false))
	tInt      = reflect.TypeOf(int(0))
	tUint     = reflect.TypeOf(uint(0))
	tUint16   = reflect.TypeOf(uint16(0))
	tString   = reflect.TypeOf(string(""))
	tDuration = reflect.TypeOf(time.Duration(0))
	pUDPAddr  = reflect.TypeOf(&net.UDPAddr{})
)

func Marshal(m interface{}) ([]byte, error) {
	v := reflect.ValueOf(m)

	if v.Type().Kind() == reflect.Ptr {
		return marshal(v.Elem())
	} else {
		return marshal(reflect.Indirect(v))
	}
}

func marshal(s reflect.Value) ([]byte, error) {
	var c strings.Builder

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			tag := t.Tag.Get("conf")

			// Marshall with MarshalConf{} interface

			if m, ok := f.Interface().(Marshaler); ok {
				// If f is a pointer type and the value is nil skips this field, leaving the buffer 'as is'
				// i.e. 'omitempty' is the default implementation
				if f.Kind() != reflect.Ptr || !f.IsNil() {
					if b, err := m.MarshalConf(); err == nil {
						fmt.Fprintf(&c, "%s = %s\n", tag, string(b))
					}
				}

				continue
			}

			// Marshal embedded structs

			if f.Kind() == reflect.Struct {
				if v, err := marshal(f); err != nil {
					return []byte(c.String()), err
				} else {
					entries := strings.Split(string(v), "\n")
					for _, e := range entries {
						if e != "" {
							fmt.Fprintf(&c, "%s.%s\n", tag, e)
						}
					}
				}

				continue
			}

			// Unmarshal built-in types

			// Marshal built-in types

			switch t.Type {
			case tBool:
				fmt.Fprintf(&c, "%s = %v\n", tag, f)

			case tInt:
				fmt.Fprintf(&c, "%s = %v\n", tag, f)

			case tUint:
				fmt.Fprintf(&c, "%s = %v\n", tag, f)

			case tString:
				fmt.Fprintf(&c, "%s = %v\n", tag, f)

			case pUDPAddr:
				fmt.Fprintf(&c, "%s = %v\n", tag, f)

			default:
				panic(errors.New(fmt.Sprintf("Cannot marshal field with type '%v'", t.Type)))
			}
		}
	}

	return []byte(c.String()), nil
}

func Unmarshal(b []byte, m interface{}) error {
	v := reflect.ValueOf(m)
	s := v.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot unmarshal %s: expected 'struct'", s.Kind())
	}

	values, err := parse(bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	return unmarshal(s, "", values)
}

func parse(r io.Reader) (map[string]string, error) {
	re := regexp.MustCompile(`^\s*(.*?)\s*=\s*(.*)\s*$`)
	m := make(map[string]string)
	s := bufio.NewScanner(r)

	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) > 0 {
			m[match[1]] = match[2]
		}
	}

	return m, s.Err()
}

func unmarshal(s reflect.Value, prefix string, values map[string]string) error {
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot unmarshal %s: expected 'struct'", s.Kind())
	}

	N := s.NumField()
	for i := 0; i < N; i++ {
		f := s.Field(i)
		t := s.Type().Field(i)
		tag := strings.TrimSpace(t.Tag.Get("conf"))

		if tag == "" || !f.CanSet() {
			continue
		}

		tag = prefix + tag

		// Unmarshal value fields with UnmarshalConf{} interface
		if u, ok := f.Addr().Interface().(Unmarshaler); ok {
			p, err := u.UnmarshalConf(tag, values)
			if err != nil {
				return err
			}

			f.Set(reflect.Indirect(reflect.ValueOf(p)))
			continue
		}

		// Unmarshal pointer fields with UnmarshalConf{} interface
		if u, ok := f.Interface().(Unmarshaler); ok {
			p, err := u.UnmarshalConf(tag, values)
			if err != nil {
				return err
			}

			f.Set(reflect.ValueOf(p))
			continue
		}

		// Unmarshal embedded structs

		if f.Kind() == reflect.Struct {
			unmarshal(f, tag+".", values)
			continue
		}

		// Unmarshal built-in types

		switch t.Type {
		case tBool:
			if value, ok := values[tag]; ok {
				if value == "true" {
					f.SetBool(true)
				} else if value == "false" {
					f.SetBool(false)
				} else {
					return fmt.Errorf("Invalid boolean value: %s:", value)
				}
			}

		case tInt:
			if value, ok := values[tag]; ok {
				i, err := strconv.ParseInt(value, 10, 0)
				if err != nil {
					return err
				}
				f.SetInt(i)
			}

		case tUint:
			if value, ok := values[tag]; ok {
				i, err := strconv.ParseUint(value, 10, 0)
				if err != nil {
					return err
				}
				f.SetUint(i)
			}

		case tUint16:
			if value, ok := values[tag]; ok {
				i, err := strconv.ParseUint(value, 10, 16)
				if err != nil {
					return err
				}
				f.SetUint(i)
			}

		case tString:
			if value, ok := values[tag]; ok {
				f.SetString(value)
			}

		case tDuration:
			if value, ok := values[tag]; ok {
				d, err := time.ParseDuration(value)
				if err != nil {
					return err
				}
				f.SetInt(int64(d))
			}

		case pUDPAddr:
			if value, ok := values[tag]; ok {
				address, err := net.ResolveUDPAddr("udp", value)
				if err != nil {
					return err
				}

				addr := net.UDPAddr{
					IP:   make(net.IP, net.IPv4len),
					Port: address.Port,
					Zone: "",
				}

				copy(addr.IP, address.IP.To4())

				f.Set(reflect.ValueOf(&addr))
			}
		}
	}

	return nil
}
