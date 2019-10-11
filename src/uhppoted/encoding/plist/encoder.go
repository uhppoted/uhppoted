package plist

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Encoder struct {
	writer io.Writer
}

type encoder func(*xml.Encoder, reflect.Value) error

var encoders = map[reflect.Type]encoder{
	reflect.TypeOf(string("")):  encodeString,
	reflect.TypeOf(bool(false)): encodeBool,
	reflect.TypeOf(int(0)):      encodeInt,
	reflect.TypeOf([]string{}):  encodeStringArray,
}

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

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		writer: w,
	}
}

func (e *Encoder) Encode(p interface{}) error {
	v := reflect.ValueOf(p)

	if v.Type().Kind() == reflect.Ptr {
		return e.encode(v.Elem())
	} else {
		return e.encode(reflect.Indirect(v))
	}
}

func (e *Encoder) encode(s reflect.Value) error {
	encoder := xml.NewEncoder(e.writer)
	encoder.Indent("", "  ")

	for _, token := range header {
		if err := encoder.EncodeToken(token); err != nil {
			return err
		}
	}

	if err := encoder.EncodeToken(body); err != nil {
		return err
	}

	if s.Kind() == reflect.Struct {
		if err := encoder.EncodeToken(dict); err != nil {
			return err
		}

		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			k := t.Name

			if err := encoder.EncodeElement(k, key); err != nil {
				return err
			}

			if g, ok := encoders[t.Type]; ok {
				if err := g(encoder, f); err != nil {
					return err
				}
			} else {
				panic(errors.New(fmt.Sprintf("Cannot encode plist field '%s' with type '%v'", k, t.Type)))
			}
		}

		if err := encoder.EncodeToken(dict.End()); err != nil {
			return err
		}

	} else {
		panic(errors.New(fmt.Sprintf("Expecting struct, got '%v'", s.Kind())))
	}

	if err := encoder.EncodeToken(body.End()); err != nil {
		return err
	}

	encoder.Flush()

	return nil
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
