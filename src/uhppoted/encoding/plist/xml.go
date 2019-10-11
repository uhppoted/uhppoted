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
	tag    string
	text   string
	parent *node
	child  *node
	next   *node
}

func (doc *document) parse(r io.Reader) error {
	doc.root = node{tag: "/"}
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
