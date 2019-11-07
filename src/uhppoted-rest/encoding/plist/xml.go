package plist

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type document struct {
	root node
}

type node struct {
	tag      string
	text     string
	parent   *node
	children struct {
		first *node
		last  *node
	}
	next *node
}

func parse(r io.Reader) (*document, error) {
	doc := document{node{tag: "/"}}
	decoder := xml.NewDecoder(r)
	p := &doc.root

	token, err := decoder.Token()
	for err == nil {
		switch v := token.(type) {
		case xml.StartElement:
			n := node{
				tag:    v.Name.Local,
				parent: p,
			}

			if p.children.first == nil {
				p.children.first = &n
			}

			if p.children.last != nil {
				p.children.last.next = &n
			}

			p.children.last = &n
			p = &n

		case xml.EndElement:
			p = p.parent

		case xml.CharData:
			if text := strings.TrimSpace(string(v)); text != "" {
				p.text = text
			}
		}

		token, err = decoder.Token()
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	return &doc, nil
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

		q := p.children.first
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
		if n.children.first != nil {
			print(n.children.first, depth+1)
		}
		if n.next != nil {
			print(n.next, depth)
		}
	}
}
