package plist

import (
	"testing"
)

var REF = `<?xml version="1.0" encoding="UTF-8"?>
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
  </array>
  <key>KeepAlive</key>
  <true/>
  <key>RunAtLoad</key>
  <false/>
  <key>StandardOutPath</key>
  <string>/usr/local/var/log/com.github.twystd.uhppoted.log</string>
  <key>StandardErrorPath</key>
  <string>/usr/local/var/log/com.github.twystd.uhppoted.err</string>
  </dict>
</plist>
`
var XML = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
 <dict>
  <key>Label</key>
  <string>com.github.twystd.uhppoted</string>
  <key>Program</key>
  <string>/usr/local/bin/uhppoted</string>
  <key>StandardErrorPath</key>
  <string>/usr/local/var/log/com.github.twystd.uhppoted.err</string>
  <key>StandardOutPath</key>
  <string>/usr/local/var/log/com.github.twystd.uhppoted.log</string>
  <key>WorkingDirectory</key>
  <string>/usr/local/var/uhppoted</string>
 </dict>
</plist>`

func TestEncode(t *testing.T) {
	p := map[string]interface{}{
		"Label":            "com.github.twystd.uhppoted",
		"Program":          "/usr/local/bin/uhppoted",
		"WorkingDirectory": "/usr/local/var/uhppoted",
		//		"KeepAlive":         true,
		"StandardOutPath":   "/usr/local/var/log/com.github.twystd.uhppoted.log",
		"StandardErrorPath": "/usr/local/var/log/com.github.twystd.uhppoted.err",
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

	//	t.Errorf("OOOPS")
}

func TestDecode(t *testing.T) {
}
