package plist

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

var XML = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.github.twystd.uhppoted</string>
    <key>Program</key>
    <string>/usr/local/bin/uhppoted</string>
    <key>WorkingDirectory</key>
    <string>/usr/local/var/uhppoted</string>
    <key>ProgramArguments</key>
    <array>
      <string>--debug</string>
      <string>--verbose</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>RunAtLoad</key>
    <false/>
    <key>StandardOutPath</key>
    <string>/usr/local/var/log/com.github.twystd.uhppoted.log</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/com.github.twystd.uhppoted.err</string>
    <key>Integer</key>
    <integer>6521</integer>
  </dict>
</plist>`

var XMLX = `<?xml version="1.0" encoding="UTF-8"?>
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>com.github.twystd.uhppoted</string>
  </dict>
</plist>`

func TestEncode(t *testing.T) {
	p := struct {
		Label             string
		Program           string
		WorkingDirectory  string
		ProgramArguments  []string
		KeepAlive         bool
		RunAtLoad         bool
		StandardOutPath   string
		StandardErrorPath string
		Integer           int
	}{
		Label:             "com.github.twystd.uhppoted",
		Program:           "/usr/local/bin/uhppoted",
		WorkingDirectory:  "/usr/local/var/uhppoted",
		ProgramArguments:  []string{"--debug", "--verbose"},
		KeepAlive:         true,
		RunAtLoad:         false,
		StandardOutPath:   "/usr/local/var/log/com.github.twystd.uhppoted.log",
		StandardErrorPath: "/usr/local/var/log/com.github.twystd.uhppoted.err",
		Integer:           6521,
	}

	buffer := bytes.Buffer{}
	encoder := NewEncoder(bufio.NewWriter(&buffer))
	err := encoder.Encode(p)
	bytes := buffer.Bytes()

	if err != nil {
		t.Errorf("plist.Encode returned unexpected error: %v", err)
		return
	}

	if string(bytes) != XML {
		t.Errorf("plist.Encodereturned unexpected XML: '%s'", string(bytes))
		return
	}
}

func TestDecode(t *testing.T) {
	p := struct {
		Label             string
		Program           string
		WorkingDirectory  string
		ProgramArguments  []string
		KeepAlive         bool
		RunAtLoad         bool
		StandardOutPath   string
		StandardErrorPath string
		Integer           int
	}{}

	decoder := NewDecoder(bytes.NewReader([]byte(XML)))
	err := decoder.Decode(&p)

	if err != nil {
		t.Fatalf("plist.Decode returned unexpected error: %v", err)
	}

	if p.Label != "com.github.twystd.uhppoted" {
		t.Errorf("plist.Decode returned unexpected string for 'Label' field: '%s'", p.Label)
	}

	if p.Program != "/usr/local/bin/uhppoted" {
		t.Errorf("plist.Decode returned unexpected string for 'Program' field: '%s'", p.Program)
	}

	if p.WorkingDirectory != "/usr/local/var/uhppoted" {
		t.Errorf("plist.Decode returned unexpected string for 'WorkingDirectory' field: '%s'", p.WorkingDirectory)
	}

	if !reflect.DeepEqual(p.ProgramArguments, []string{"--debug", "--verbose"}) {
		t.Errorf("plist.Decode returned unexpected string array for 'ProgramArguments' field: '%v'", p.ProgramArguments)
	}

	if !p.KeepAlive {
		t.Errorf("plist.Decode returned unexpected bool for 'KeepAlive' field: '%v'", p.KeepAlive)
	}

	if p.RunAtLoad {
		t.Errorf("plist.Decode returned unexpected bool for 'RunAtLoad' field: '%v'", p.RunAtLoad)
	}

	if p.StandardOutPath != "/usr/local/var/log/com.github.twystd.uhppoted.log" {
		t.Errorf("plist.Decode returned unexpected string for 'StandardOutPath' field: '%s'", p.StandardOutPath)
	}

	if p.StandardErrorPath != "/usr/local/var/log/com.github.twystd.uhppoted.err" {
		t.Errorf("plist.Decode returned unexpected string for 'StandardErrorPath' field: '%s'", p.StandardErrorPath)
	}

	if p.Integer != 6521 {
		t.Errorf("plist.Decode returned unexpected integer for 'Integer' field: '%v'", p.Integer)
	}
}
