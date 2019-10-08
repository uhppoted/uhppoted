package commands

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var plist = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
      <string>{{.Label}}</string>
    <key>Program</key>
      <string>{{.Program}}</string>
    <key>WorkingDirectory</key>
      <string>{{.WorkingDirectory}}</string>
    <key>ProgramArguments</key>
      <array></array>
    <key>KeepAlive</key>
      <{{.KeepAlive}}/>
    <key>RunAtLoad</key>
      <{{.RunAtLoad}}/>
    <key>StandardOutPath</key>
      <string>{{.StandardOutPath}}</string>
    <key>StandardErrorPath</key>
      <string>{{.StandardErrorPath}}</string>
  </dict>
</plist>
`

type data = struct {
	Label             string
	Program           string
	WorkingDirectory  string
	KeepAlive         bool
	RunAtLoad         bool
	StandardOutPath   string
	StandardErrorPath string
}

type Daemonize struct {
}

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	if err := c.launchd(); err != nil {
		return err
	}

	if err := c.mkdirs(); err != nil {
		return err
	}

	if err := c.firewall(); err != nil {
		return err
	}

	fmt.Println("   ... com.github.twystd.uhppoted registered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("   sudo launchctl load /Libary/LaunchDaemons/com.github.twystd.uhppoted")
	fmt.Println()

	return nil
}

func (c *Daemonize) launchd() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	d := data{
		Label:             "com.github.twystd.uhppoted",
		Program:           executable,
		WorkingDirectory:  "/usr/local/var/com.github.twystd.uhppoted",
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   "/usr/local/var/log/com.github.twystd.uhppoted.log",
		StandardErrorPath: "/usr/local/var/log/com.github.twystd.uhppoted.err",
	}

	path := filepath.Join("/Library/LaunchDaemons", "com.github.twystd.uhppoted.plist")
	_, err = os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !os.IsNotExist(err) {
		dict, err := c.parse(path)
		if err != nil {
			return err
		}

		if label, ok := dict["Label"]; ok {
			d.Label = label
		}

		if directory, ok := dict["WorkingDirectory"]; ok {
			d.WorkingDirectory = directory
		}

		if keepalive, ok := dict["KeepAlive"]; ok && keepalive == "false" {
			d.KeepAlive = false
		}

		if runatload, ok := dict["RunAtLoad"]; ok && runatload == "false" {
			d.RunAtLoad = false
		}

		if stdout, ok := dict["StandardOutPath"]; ok {
			d.StandardOutPath = stdout
		}

		if stderr, ok := dict["StandardErrorPath"]; ok {
			d.StandardErrorPath = stderr
		}
	}

	return c.daemonize(path, d)
}

// TODO rework this to parse XML plist files more robustly
func (c *Daemonize) parse(path string) (map[string]string, error) {
	dict := make(map[string]string)
	f, err := os.Open(path)
	if err != nil {
		return dict, err
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)

	text := ""
	key := ""
	value := ""

	for {
		token, err := decoder.Token()

		if err != nil {
			if err != io.EOF {
				return dict, err
			}
			break
		}

		if _, ok := token.(xml.StartElement); ok {
		} else if end, ok := token.(xml.EndElement); ok {
			switch end.Name.Local {
			case "key":
				key = text
			case "string":
				value = text
				dict[key] = value
			case "true":
				dict[key] = "true"
			case "false":
				dict[key] = "false"
			}
		} else if char, ok := token.(xml.CharData); ok {
			text = string(char)
		}
	}

	return dict, nil
}

func (c *Daemonize) daemonize(path string, d data) error {
	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	t := template.Must(template.New("com.github.twystd.uhppoted.plist").Parse(plist))
	err = t.Execute(f, d)
	if err != nil {
		return err
	}

	return nil
}

func (c *Daemonize) mkdirs() error {
	dir := "/usr/local/var/com.github.twystd.uhppoted"

	fmt.Printf("   ... creating '%s'\n", dir)

	return os.MkdirAll(dir, 0644)
}

func (c *Daemonize) firewall() error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Println("   *** WARNING: adding 'uhppoted' to the application firewall and unblocking incoming connections")
	fmt.Println("   ***")
	fmt.Println()

	path, err := os.Executable()
	if err != nil {
		fmt.Errorf("Failed to get path to executable: %v\n", err)
		return err
	}

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		fmt.Errorf("ERROR: Failed to retrieve application firewall global state (%v)\n", err)
		return err
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to disable the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to add 'uhppoted' to the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to unblock 'uhppoted' on the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to re-enable the application firewall (%v)\n", err)
			return err
		}

		fmt.Println()
	}

	return nil
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Daemonizes uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Daemonizes uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}
