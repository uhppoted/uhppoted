package plist

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type document struct {
	root node
}

type node struct {
	tag    string
	text   string
	parent *node
	child  *node
	next   *node
}

type encoder func(*xml.Encoder, reflect.Value) error
type decoder func(*node, string, reflect.Value) error

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

func Encode(p interface{}, w io.Writer) error {
	v := reflect.ValueOf(p)

	if v.Type().Kind() == reflect.Ptr {
		return encode(v.Elem(), w)
	} else {
		return encode(reflect.Indirect(v), w)
	}
}

func encode(s reflect.Value, w io.Writer) error {
	encoder := xml.NewEncoder(w)
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

func Decode(r io.Reader, p interface{}) error {
	doc := document{node{tag: "/"}}
	if err := doc.parse(r); err != nil {
		return err
	}

	doc.print()

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
	p := n.child

	for {
		if p == nil {
			break
		} else if p.tag != "string" {
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

func (doc *document) parse(r io.Reader) error {
	decoder := xml.NewDecoder(r)

	var current *node = &doc.root
	var peer *node = nil

	for {
		token, err := decoder.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if start, ok := token.(xml.StartElement); ok {
			n := node{
				tag:    start.Name.Local,
				text:   "",
				parent: nil,
				child:  nil,
				next:   nil,
			}

			if current != nil {
				n.parent = current
				if current.child == nil {
					current.child = &n
				}
			} else if peer != nil {
				n.parent = peer.parent
				peer.next = &n
			}

			current = &n
			peer = &n
			continue
		}

		if _, ok := token.(xml.EndElement); ok {
			if current == nil {
				if peer != nil {
					peer = peer.parent
				} else {
					peer = nil
				}
			}
			current = nil
			continue
		}

		if chardata, ok := token.(xml.CharData); ok {
			if text := strings.TrimSpace(string(chardata)); text != "" {
				if current != nil {
					current.text = text
				}
			}
			continue
		}
	}

	return nil
}

func (doc *document) query(xpath string) (*node, error) {
	re := regexp.MustCompile(`(\w+)(?:\[text="(.*?)"\])?`)
	segments := strings.Split(xpath, "/")

	var p *node = &doc.root

	for _, s := range segments[1:] {
		if s == "next()" {
			p = p.next
			continue
		}

		match := re.FindStringSubmatch(s)
		if match == nil {
			return nil, fmt.Errorf("Invalid query: '%s'", xpath)
		}

		tag := match[1]
		text := match[2]

		q := p.child
		for {
			if q == nil {
				return nil, nil
			}

			if q.tag == tag && (text == "" || q.text == text) {
				p = q
				break
			} else {
				q = q.next
			}
		}
	}

	return p, nil
}

func (doc *document) print() {
	print(&doc.root, 0)
}

func print(n *node, depth int) {
	indent := strings.Repeat(" ", depth)
	if n != nil {
		fmt.Printf("%p%s (%v)\n", n, indent, *n)
		if n.child != nil {
			print(n.child, depth+1)
		}
		if n.next != nil {
			print(n.next, depth)
		}
	}
}
