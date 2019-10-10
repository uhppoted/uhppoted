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

var (
	tString      = reflect.TypeOf(string(""))
	tBool        = reflect.TypeOf(bool(false))
	tInteger     = reflect.TypeOf(int(0))
	tStringArray = reflect.TypeOf([]string{})
)

var instruction = xml.ProcInst{"xml", []byte(`version="1.0" encoding="UTF-8"`)}
var doctype = xml.Directive([]byte(`DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))
var newline = xml.CharData("\n")
var dict = xml.StartElement{Name: xml.Name{"", "dict"}}
var key = xml.StartElement{Name: xml.Name{"", "key"}}
var integer = xml.StartElement{Name: xml.Name{"", "integer"}}
var array = xml.StartElement{Name: xml.Name{"", "array"}}
var boolean = map[bool]xml.StartElement{
	true:  xml.StartElement{Name: xml.Name{"", "true"}},
	false: xml.StartElement{Name: xml.Name{"", "false"}},
}

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

			switch t.Type {
			case tString:
				if err := encoder.Encode(f.String()); err != nil {
					return buffer.Bytes(), err
				}

			case tBool:
				v := xml.CharData("")
				element := boolean[f.Bool()]
				if err := encoder.EncodeElement(v, element); err != nil {
					return buffer.Bytes(), err
				}

			case tInteger:
				v := xml.CharData(strconv.FormatInt(f.Int(), 10))
				if err := encoder.EncodeElement(v, integer); err != nil {
					return buffer.Bytes(), err
				}

			case tStringArray:
				if err := encoder.EncodeToken(array); err != nil {
					return buffer.Bytes(), err
				}

				for j := 0; j < f.Len(); j++ {
					if err := encoder.Encode(f.Index(j).String()); err != nil {
						return buffer.Bytes(), err
					}
				}

				if err := encoder.EncodeToken(array.End()); err != nil {
					return buffer.Bytes(), err
				}

			default:
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

func Decode(bytes []byte) (interface{}, error) {
	return nil, nil
}
