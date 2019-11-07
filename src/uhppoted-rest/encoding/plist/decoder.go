package plist

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Decoder struct {
	reader io.Reader
}

type decoder func(*node, string, reflect.Value) error

var decoders = map[reflect.Type]decoder{
	reflect.TypeOf(string("")):  decodeString,
	reflect.TypeOf(bool(false)): decodeBool,
	reflect.TypeOf(int(0)):      decodeInt,
	reflect.TypeOf([]string{}):  decodeStringArray,
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader: r,
	}
}

func (d *Decoder) Decode(p interface{}) error {
	doc, err := parse(d.reader)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(p)
	s := v.Elem()

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			k := t.Name
			q := fmt.Sprintf(`/plist/dict/key[text="%s"]/next()`, k)

			n, err := doc.query(q)
			if err != nil {
				return err
			}

			if n != nil {
				if !f.CanSet() {
					return fmt.Errorf("Cannot set struct field '%s'", k)
				}

				if g, ok := decoders[t.Type]; ok {
					if err := g(n, k, f); err != nil {
						return err
					}
				} else {
					panic(errors.New(fmt.Sprintf("Cannot decode  plist field '%s' with type '%v'", k, t.Type)))
				}
			}
		}
	} else {
		panic(errors.New(fmt.Sprintf("Expecting struct, got '%v'", s.Kind())))
	}

	return nil
}

func decodeString(n *node, field string, f reflect.Value) error {
	if n.tag != "string" {
		return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'string'", n.tag, field)
	}

	f.SetString(n.text)
	return nil
}

func decodeBool(n *node, field string, f reflect.Value) error {
	if n.tag == "true" {
		f.SetBool(true)
	} else if n.tag == "false" {
		f.SetBool(false)
	} else {
		return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'bool'", n.tag, field)
	}

	return nil
}

func decodeInt(n *node, field string, f reflect.Value) error {
	if n.tag != "integer" {
		return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'integer'", n.tag, field)
	}

	if ivalue, err := strconv.ParseInt(n.text, 10, 64); err != nil {
		return err
	} else {
		f.SetInt(ivalue)
	}

	return nil
}

func decodeStringArray(n *node, field string, f reflect.Value) error {
	if n.tag != "array" {
		return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'array'", n.tag, field)
	}

	strings := []string{}

	p := n.children.first
	for p != nil {
		if p.tag != "string" {
			return fmt.Errorf("Invalid plist XML array element '%s' for field '%s': expected 'string'", p.tag, field)
		}

		strings = append(strings, p.text)
		p = p.next
	}

	values := reflect.MakeSlice(reflect.Indirect(f).Type(), len(strings), len(strings))
	for i, s := range strings {
		values.Index(i).SetString(s)
	}

	f.Set(values)

	return nil
}
