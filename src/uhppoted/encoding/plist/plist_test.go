package plist

import (
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
    <true></true>
    <key>RunAtLoad</key>
    <false></false>
    <key>StandardOutPath</key>
    <string>/usr/local/var/log/com.github.twystd.uhppoted.log</string>
    <key>StandardErrorPath</key>
    <string>/usr/local/var/log/com.github.twystd.uhppoted.err</string>
    <key>Integer</key>
    <integer>6521</integer>
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

	bytes, err := Encode(p)

	if err != nil {
		t.Errorf("plist.Encode(%v) returned unexpected error: %v", p, err)
		return
	}

	if string(bytes) != XML {
		t.Errorf("plist.Encode(%v) returned invalid XML: '%s'", p, string(bytes))
		return
	}
}

func TestDecode(t *testing.T) {
}
