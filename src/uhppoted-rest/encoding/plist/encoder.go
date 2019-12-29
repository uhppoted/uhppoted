package plist

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type Encoder struct {
	writer  io.Writer
	encoder *xml.Encoder
}

type encoder func(*Encoder, reflect.Value) error

var encoders = map[reflect.Type]encoder{
	reflect.TypeOf(string("")):  encodeString,
	reflect.TypeOf(bool(false)): encodeBool,
	reflect.TypeOf(int(0)):      encodeInt,
	reflect.TypeOf([]string{}):  encodeStringArray,
}

var instruction = xml.ProcInst{
	Target: "xml",
	Inst:   []byte(`version="1.0" encoding="UTF-8"`),
}

var doctype = xml.Directive([]byte(`DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))
var newline = xml.CharData("\n")
var dict = xml.StartElement{
	Name: xml.Name{
		Space: "",
		Local: "dict",
	},
}

var key = xml.StartElement{
	Name: xml.Name{
		Space: "",
		Local: "key"},
}

var header = []xml.Token{
	instruction,
	newline,
	doctype,
	newline,
}

var body = xml.StartElement{
	Name: xml.Name{
		Space: "",
		Local: "plist",
	},
	Attr: []xml.Attr{
		xml.Attr{
			Name: xml.Name{
				Space: "",
				Local: "version",
			},
			Value: "1.0"},
	},
}

func NewEncoder(w io.Writer) *Encoder {
	e := xml.NewEncoder(w)
	e.Indent("", "  ")

	return &Encoder{
		writer:  w,
		encoder: e,
	}
}

func (e *Encoder) Encode(p interface{}) error {
	v := reflect.ValueOf(p)

	if v.Type().Kind() == reflect.Ptr {
		return e.encode(v.Elem())
	}

	return e.encode(reflect.Indirect(v))
}

func (e *Encoder) encode(s reflect.Value) error {

	for _, token := range header {
		if err := e.encoder.EncodeToken(token); err != nil {
			return err
		}
	}

	if err := e.encoder.EncodeToken(body); err != nil {
		return err
	}

	if s.Kind() == reflect.Struct {
		if err := e.encoder.EncodeToken(dict); err != nil {
			return err
		}

		N := s.NumField()

		for i := 0; i < N; i++ {
			f := s.Field(i)
			t := s.Type().Field(i)
			k := t.Name

			if err := e.encoder.EncodeElement(k, key); err != nil {
				return err
			}

			if g, ok := encoders[t.Type]; ok {
				if err := g(e, f); err != nil {
					return err
				}
			} else {
				panic(fmt.Errorf("Cannot encode plist field '%s' with type '%v'", k, t.Type))
			}
		}

		if err := e.encoder.EncodeToken(dict.End()); err != nil {
			return err
		}

	} else {
		panic(fmt.Errorf("Expecting struct, got '%v'", s.Kind()))
	}

	if err := e.encoder.EncodeToken(body.End()); err != nil {
		return err
	}

	e.encoder.Flush()

	return nil
}

func encodeString(e *Encoder, f reflect.Value) error {
	value := f.String()
	element := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "string",
		},
	}

	return e.encoder.EncodeElement(value, element)
}

// Aaaaaargh! MacOS requires the <true/> and <false/> form while an unmodified
// Go XML encoder can only write <true></true>. In 2019!!
func encodeBool(e *Encoder, f reflect.Value) error {
	v := fmt.Sprintf("\n    <%v/>", f.Bool())
	_, err := e.writer.Write([]byte(v))

	return err
}

func encodeInt(e *Encoder, f reflect.Value) error {
	value := xml.CharData(strconv.FormatInt(f.Int(), 10))
	element := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "integer",
		},
	}

	return e.encoder.EncodeElement(value, element)
}

func encodeStringArray(e *Encoder, f reflect.Value) error {
	start := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "array",
		},
	}
	end := start.End()

	if err := e.encoder.EncodeToken(start); err != nil {
		return err
	}

	for j := 0; j < f.Len(); j++ {
		if err := encodeString(e, f.Index(j)); err != nil {
			return err
		}
	}

	if err := e.encoder.EncodeToken(end); err != nil {
		return err
	}

	return nil
}
