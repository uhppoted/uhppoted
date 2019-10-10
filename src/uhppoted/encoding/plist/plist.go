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

type encoding func(*xml.Encoder, reflect.Value) error

var dispatch = map[reflect.Type]encoding{
	reflect.TypeOf(string("")):  encodeString,
	reflect.TypeOf(bool(false)): encodeBool,
	reflect.TypeOf(int(0)):      encodeInt,
	reflect.TypeOf([]string{}):  encodeStringArray,
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

			if g, ok := dispatch[t.Type]; ok {
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

func Decode(bytes []byte) (interface{}, error) {
	return nil, nil
}
