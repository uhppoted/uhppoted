package commands

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Undaemonize struct {
}

func (c *Undaemonize) Execute(ctx Context) error {
	fmt.Println("   ... undaemonizing")

	if err := c.launchd(); err != nil {
		return err
	}

	if err := c.rmdirs(); err != nil {
		return err
	}

	if err := c.firewall(); err != nil {
		return err
	}

	fmt.Println("   ... com.github.twystd.uhppoted unregistered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   Any uhppoted log files can still be found in directory /usr/local/var/log and should be removed manually")
	fmt.Println()

	return nil
}

func (c *Undaemonize) launchd() error {
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

	if os.IsNotExist(err) {
		fmt.Printf("   ... %s does not exist - nothing to do for launchd\n", path)
		return nil
	}

	dict, err := c.parse(path)
	if err != nil {
		return err
	}

	if label, ok := dict["Label"]; ok {
		d.Label = label
	}

	if program, ok := dict["Program"]; ok {
		d.Program = program
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

	return c.undaemonize(path, d)
}

// TODO rework this to parse XML plist files more robustly
func (c *Undaemonize) parse(path string) (map[string]string, error) {
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

func (c *Undaemonize) undaemonize(path string, d data) error {
	fmt.Printf("   ... unloading LaunchDaemon\n")
	cmd := exec.Command("launchctl", "unload", path)
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		fmt.Errorf("ERROR: Failed to unload '%s' (%v)\n", d.Label, err)
		return err
	}

	fmt.Printf("   ... removing '%s'\n", path)
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Undaemonize) rmdirs() error {
	dir := "/usr/local/var/com.github.twystd.uhppoted"

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}

func (c *Undaemonize) firewall() error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Println("   *** WARNING: removing 'uhppoted' to the application firewall")
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

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--remove", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to remove 'uhppoted' from the application firewall (%v)\n", err)
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

func (c *Undaemonize) Cmd() string {
	return "undaemonize"
}

func (c *Undaemonize) Description() string {
	return "Deregisters uhppoted as a service/daemon"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println("Usage: uhppoted undaemonize")
	fmt.Println()
	fmt.Println(" Deregisters uhppoted from launchd as a service/daemon")
	fmt.Println()
}
