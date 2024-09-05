package main

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
)

var funcs = template.FuncMap{}

func generate() {
	fsys := templates
	dir := filepath.Join("generated", "uhppoted-python")

	if err := os.MkdirAll(dir, 0750); err != nil {
		errorf("uhppoted-python", "%v", err)
	} else if t, err := template.New("uhppoted").Funcs(funcs).ParseFS(fsys, "templates/uhppoted-python/*"); err != nil {
		errorf("uhppoted-python", "%v", err)
	} else {
		list := t.Templates()
		data := map[string]any{}
		for _, v := range list {
			infof("uhppoted-python", "... processing template %v", v.Name())
			var b bytes.Buffer

			if err := v.Execute(&b, data); err != nil {
				errorf("uhppoted-python", "%v", err)
			} else {
				path := filepath.Join(dir, v.Name())
				if file, err := os.Create(path); err != nil {
					errorf("uhppoted-python", "%v", err)
				} else {
					if _, err := file.Write(b.Bytes()); err != nil {
						errorf("uhppoted-python", "%v", err)
					}
					file.Close()
				}
			}
		}
	}
}
