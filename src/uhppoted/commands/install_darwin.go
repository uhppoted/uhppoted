package commands

import (
	"fmt"
	"log"
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
      <array>
      </array>
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

type Install struct {
}

func (c *Install) Execute(ctx Context) error {
	fmt.Println("...... installing")

	if err := c.launchd(); err != nil {
		return err
	}

	if err := c.mkdirs(); err != nil {
		return err
	}

	if err := c.firewall(); err != nil {
		return err
	}

	return nil
}

func (c *Install) launchd() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	data := struct {
		Label             string
		Program           string
		WorkingDirectory  string
		KeepAlive         bool
		RunAtLoad         bool
		StandardOutPath   string
		StandardErrorPath string
	}{
		Label:             "com.github.twystd.uhppoted",
		Program:           executable,
		WorkingDirectory:  "/usr/local/var/com.github.twystd.uhppoted",
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   "/usr/local/var/log/com.github.twystd.uhppoted.log",
		StandardErrorPath: "/usr/local/var/log/com.github.twystd.uhppoted.err",
	}

	path := filepath.Join("/Library/LaunchDaemons", data.Label+".plist")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	t := template.Must(template.New("com.github.twystd.uhppoted.plist").Parse(plist))
	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Install) mkdirs() error {
	if err := os.MkdirAll("/usr/local/var/com.github.twystd.uhppoted", 0644); err != nil {
		return err
	}

	return nil
}

func (c *Install) firewall() error {
	log.Println()
	log.Println("   ***")
	log.Println("   *** WARNING: adding 'uhppoted' to the application firewall and unblocking incoming connections")
	log.Println("   ***")
	log.Println()

	path, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get path to executable: %v\n", err)
		return err
	}

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	log.Printf("   > %s", out)
	if err != nil {
		log.Fatalf("ERROR: Failed to retrieve application firewall global state (%v)\n", err)
		return err
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to disable the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to add 'uhppoted' to the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to unblock 'uhppoted' on the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to re-enable the application firewall (%v)\n", err)
			return err
		}

		log.Println()
	}

	return nil
}

func (c *Install) CLI() string {
	return "install"
}

func (c *Install) Description() string {
	return "Installs uhppoted as a service/daemon"
}

func (c *Install) Usage() string {
	return ""
}

func (c *Install) Help() {
	fmt.Println("Usage: uhppoted install")
	fmt.Println()
	fmt.Println(" Installs uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}
