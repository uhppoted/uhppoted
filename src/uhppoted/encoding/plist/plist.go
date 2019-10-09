package plist

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"sort"
)

type PListItem struct {
	Key   string
	Value interface{}
}

func Encode(p map[string]interface{}) ([]byte, error) {
	var buffer bytes.Buffer

	instruction := xml.ProcInst{"xml", []byte(`version="1.0" encoding="UTF-8"`)}
	doctype := xml.Directive([]byte(`DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"`))
	newline := xml.CharData("\n")

	body := xml.StartElement{
		Name: xml.Name{"", "plist"},
		Attr: []xml.Attr{
			xml.Attr{xml.Name{"", "version"}, "1.0"}},
	}

	dict := xml.StartElement{
		Name: xml.Name{"", "dict"},
		Attr: []xml.Attr{},
	}

	key := xml.StartElement{
		Name: xml.Name{"", "key"},
		Attr: []xml.Attr{},
	}

	header := []xml.Token{
		instruction,
		newline,
		doctype,
		newline,
	}

	encoder := xml.NewEncoder(bufio.NewWriter(&buffer))
	encoder.Indent("", " ")

	for _, token := range header {
		if err := encoder.EncodeToken(token); err != nil {
			return buffer.Bytes(), err
		}
	}

	if err := encoder.EncodeToken(body); err != nil {
		return buffer.Bytes(), err
	}

	if err := encoder.EncodeToken(dict); err != nil {
		return buffer.Bytes(), err
	}

	// Sort keys just to make testing a bit saner
	var keys []string
	for k := range p {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if err := encoder.EncodeElement(k, key); err != nil {
			return buffer.Bytes(), err
		}
		if err := encoder.Encode(p[k]); err != nil {
			return buffer.Bytes(), err
		}
	}

	if err := encoder.EncodeToken(dict.End()); err != nil {
		return buffer.Bytes(), err
	}

	if err := encoder.EncodeToken(body.End()); err != nil {
		return buffer.Bytes(), err
	}

	encoder.Flush()

	fmt.Println("----")
	fmt.Println(string(buffer.Bytes()))
	fmt.Println("---")

	return buffer.Bytes(), nil
}

func Decode(bytes []byte) ([]PListItem, error) {
	return []PListItem{}, nil
}
