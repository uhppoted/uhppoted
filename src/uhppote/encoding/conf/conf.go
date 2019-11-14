package conf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"reflect"
	"regexp"
	"strings"
)

type Unmarshaler interface {
	UnmarshalConf(string) (interface{}, error)
}

var (
	pUDPAddr = reflect.TypeOf(&net.UDPAddr{})
)

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

	N := s.NumField()
	for i := 0; i < N; i++ {
		f := s.Field(i)
		t := s.Type().Field(i)
		tag := strings.TrimSpace(t.Tag.Get("conf"))
		value, found := values[tag]

		if tag == "" || !found || !f.CanSet() {
			continue
		}

		// Unmarshall value fields with UnmarshalConf{} interface
		if u, ok := f.Addr().Interface().(Unmarshaler); ok {
			p, err := u.UnmarshalConf(value)
			if err != nil {
				return err
			}

			f.Set(reflect.Indirect(reflect.ValueOf(p)))

			continue
		}

		// Unmarshall pointer fields with UnmarshalConf{} interface
		if u, ok := f.Interface().(Unmarshaler); ok {
			p, err := u.UnmarshalConf(value)
			if err != nil {
				return err
			}

			f.Set(reflect.ValueOf(p))

			continue
		}

		// Unmarshal built-in types
		switch t.Type {
		case pUDPAddr:
			address, err := net.ResolveUDPAddr("udp", value)
			if err != nil {
				return err
			}

			f.Set(reflect.ValueOf(address))
		}
	}

	return nil
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
