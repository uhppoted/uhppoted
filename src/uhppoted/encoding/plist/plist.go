package plist

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var instruction = xml.ProcInst{"xml", []byte(`version="1.0" encoding="UTF-8"`)}
var doctype = xml.Directive([]byte(`DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))
var newline = xml.CharData("\n")
var dict = xml.StartElement{Name: xml.Name{"", "dict"}}
var key = xml.StartElement{Name: xml.Name{"", "key"}}

var header = []xml.Token{
	instruction,
	newline,
	doctype,
	newline,
}

var body = xml.StartElement{
	Name: xml.Name{"", "plist"},
	Attr: []xml.Attr{
		xml.Attr{xml.Name{"", "version"}, "1.0"}},
}

type encoder func(*xml.Encoder, reflect.Value) error
type decoder func(*xml.Decoder, string, reflect.Value) error

var encoders = map[reflect.Type]encoder{
	reflect.TypeOf(string("")):  encodeString,
	reflect.TypeOf(bool(false)): encodeBool,
	reflect.TypeOf(int(0)):      encodeInt,
	reflect.TypeOf([]string{}):  encodeStringArray,
}

var decoders = map[reflect.Type]decoder{
	reflect.TypeOf(string("")):  decodeString,
	reflect.TypeOf(bool(false)): decodeBool,
	reflect.TypeOf(int(0)):      decodeInt,
	reflect.TypeOf([]string{}):  decodeStringArray,
}

func Encode(p interface{}) ([]byte, error) {
	v := reflect.ValueOf(p)

	if v.Type().Kind() == reflect.Ptr {
		return encode(v.Elem())
	} else {
		return encode(reflect.Indirect(v))
	}
}

func encode(s reflect.Value) ([]byte, error) {
	buffer := bytes.Buffer{}
	encoder := xml.NewEncoder(bufio.NewWriter(&buffer))
	encoder.Indent("", "  ")

	for _, token := range header {
		if err := encoder.EncodeToken(token); err != nil {
			return buffer.Bytes(), err
		}
	}

	if err := encoder.EncodeToken(body); err != nil {
		return buffer.Bytes(), err
	}

	if s.Kind() == reflect.Struct {
		if err := encoder.EncodeToken(dict); err != nil {
			return buffer.Bytes(), err
		}

		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			k := t.Name

			if err := encoder.EncodeElement(k, key); err != nil {
				return buffer.Bytes(), err
			}

			if g, ok := encoders[t.Type]; ok {
				if err := g(encoder, f); err != nil {
					return buffer.Bytes(), err
				}
			} else {
				panic(errors.New(fmt.Sprintf("Cannot encode plist field with type '%v'", t.Type)))
			}
		}

		if err := encoder.EncodeToken(dict.End()); err != nil {
			return buffer.Bytes(), err
		}

	} else {
		panic(errors.New(fmt.Sprintf("Expecting struct, got '%v'", s.Kind())))
	}

	if err := encoder.EncodeToken(body.End()); err != nil {
		return buffer.Bytes(), err
	}

	encoder.Flush()

	return buffer.Bytes(), nil
}

func encodeString(e *xml.Encoder, f reflect.Value) error {
	value := f.String()
	element := xml.StartElement{Name: xml.Name{"", "string"}}

	return e.EncodeElement(value, element)
}

func encodeBool(e *xml.Encoder, f reflect.Value) error {
	value := xml.CharData("")
	element := xml.StartElement{Name: xml.Name{"", strconv.FormatBool(f.Bool())}}

	return e.EncodeElement(value, element)
}

func encodeInt(e *xml.Encoder, f reflect.Value) error {
	value := xml.CharData(strconv.FormatInt(f.Int(), 10))
	element := xml.StartElement{Name: xml.Name{"", "integer"}}

	return e.EncodeElement(value, element)
}

func encodeStringArray(e *xml.Encoder, f reflect.Value) error {
	start := xml.StartElement{Name: xml.Name{"", "array"}}
	end := start.End()

	if err := e.EncodeToken(start); err != nil {
		return err
	}

	for j := 0; j < f.Len(); j++ {
		if err := encodeString(e, f.Index(j)); err != nil {
			return err
		}
	}

	if err := e.EncodeToken(end); err != nil {
		return err
	}

	return nil
}

func Decode(b []byte, p interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(b))

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		element, ok := token.(xml.StartElement)
		if ok {
			if element.Name.Local == "plist" {
				break
			}
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "dict" {
				break
			}
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.EndElement); ok {
			if element.Name.Local == "dict" {
				break
			}
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "key" {
				key, err := decodeKey(decoder)
				if err != nil {
					return err
				}

				err = decodeValue(decoder, key, p)
				if err != nil {
					return err
				}
			}
		}
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.EndElement); ok {
			if element.Name.Local == "plist" {
				break
			}
		}
	}

	return nil
}

func decodeKey(d *xml.Decoder) (string, error) {
	token, err := d.Token()
	if err != nil {
		return "", err
	}

	chardata, ok := token.(xml.CharData)
	if !ok {
		return "", errors.New("Invalid plist 'key'")
	}

	key := string(chardata)

	token, err = d.Token()
	_, ok = token.(xml.EndElement)
	if !ok {
		return "", errors.New("Missing plist </key> element")
	}

	return key, nil
}

func decodeValue(d *xml.Decoder, key string, p interface{}) error {
	v := reflect.ValueOf(p)
	s := v.Elem()

	if s.Kind() == reflect.Struct {
		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			k := t.Name

			if k == key {
				if g, ok := decoders[t.Type]; ok {
					if err := g(d, k, f); err != nil {
						return err
					}
				} else {
					panic(errors.New(fmt.Sprintf("Cannot decode  plist field with type '%v'", t.Type)))
				}
			}
		}

	} else {
		panic(errors.New(fmt.Sprintf("Expecting struct, got '%v'", s.Kind())))
	}

	return nil
}

func decodeString(d *xml.Decoder, name string, f reflect.Value) error {
	if !f.CanSet() {
		return fmt.Errorf("Cannot set struct field '%s'", name)
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "string" {
				break
			}

			return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'string'", element.Name.Local, name)
		}
	}

	token, err := d.Token()
	if err != nil {
		return err
	}

	chardata, ok := token.(xml.CharData)
	if !ok {
		return fmt.Errorf("Invalid plist 'string' for '%s'", name)
	}

	f.SetString(string(chardata))

	return nil
}

func decodeBool(d *xml.Decoder, name string, f reflect.Value) error {
	if !f.CanSet() {
		return fmt.Errorf("Cannot set struct field '%s'", name)
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "true" {
				f.SetBool(true)
				return nil
			}

			if element.Name.Local == "false" {
				f.SetBool(false)
				return nil
			}

			return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'string'", element.Name.Local, name)
		}
	}

	return fmt.Errorf("Missing plist value for field '%s'", name)
}

func decodeInt(d *xml.Decoder, name string, f reflect.Value) error {
	if !f.CanSet() {
		return fmt.Errorf("Cannot set struct field '%s'", name)
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "integer" {
				break
			}

			return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'integer'", element.Name.Local, name)
		}
	}

	token, err := d.Token()
	if err != nil {
		return err
	}

	chardata, ok := token.(xml.CharData)
	if !ok {
		return fmt.Errorf("Invalid plist 'integer' for '%s'", name)
	}

	ivalue, err := strconv.ParseInt(string(chardata), 10, 64)
	if err != nil {
		return err
	}

	f.SetInt(ivalue)

	return nil
}

func decodeStringArray(d *xml.Decoder, name string, f reflect.Value) error {
	if !f.CanSet() {
		return fmt.Errorf("Cannot set struct field '%s'", name)
	}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local == "array" {
				break
			}

			return fmt.Errorf("Invalid plist XML element '%s' for field '%s': expected 'array'", element.Name.Local, name)
		}
	}

	strings := []string{}

	for {
		token, err := d.Token()
		if err != nil {
			return err
		}

		if element, ok := token.(xml.EndElement); ok {
			if element.Name.Local == "array" {
				break
			}
		}

		if element, ok := token.(xml.StartElement); ok {
			if element.Name.Local != "string" {
				return fmt.Errorf("Invalid plist XML array element '%s' for field '%s': expected 'string'", element.Name.Local, name)
			}

			token, err := d.Token()
			if err != nil {
				return err
			}

			chardata, ok := token.(xml.CharData)
			if !ok {
				return fmt.Errorf("Invalid plist 'string' for '%s'", name)
			}

			strings = append(strings, string(chardata))
		}
	}

	values := reflect.MakeSlice(reflect.Indirect(f).Type(), len(strings), len(strings))

	for i, s := range strings {
		values.Index(i).SetString(s)
	}

	f.Set(values)

	return nil
}
