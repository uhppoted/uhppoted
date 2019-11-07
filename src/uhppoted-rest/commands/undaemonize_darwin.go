package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	xml "uhppoted-rest/encoding/plist"
)

type Undaemonize struct {
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{}
}

func (c *Undaemonize) Name() string {
	return "undaemonize"
}

func (c *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (c *Undaemonize) Description() string {
	return "Deregisters uhppoted-rest as a service/daemon"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest undaemonize")
	fmt.Println()
	fmt.Println("    Deregisters uhppoted-rest from launchd as a service/daemon")
	fmt.Println()
}

func (c *Undaemonize) Execute(ctx Context) error {
	fmt.Println("   ... undaemonizing")

	path := filepath.Join("/Library/LaunchDaemons", "com.github.twystd.uhppoted-rest.plist")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... %s does not exist - nothing to do for launchd\n", path)
		return nil
	}

	p, err := c.parse(path)
	if err != nil {
		return err
	}

	if err := c.launchd(path, *p); err != nil {
		return err
	}

	if err := c.logrotate(); err != nil {
		return err
	}

	if err := c.rmdirs(*p); err != nil {
		return err
	}

	if err := c.firewall(*p); err != nil {
		return err
	}

	fmt.Println("   ... com.github.twystd.uhppoted-rest unregistered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   Any uhppoted-rest log files can still be found in directory /usr/local/var/log and should be removed manually")
	fmt.Println()

	return nil
}

func (c *Undaemonize) parse(path string) (*info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := info{}
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *Undaemonize) launchd(path string, d info) error {
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

func (c *Undaemonize) logrotate() error {
	path := filepath.Join("/etc/newsyslog.d", "uhppoted-rest.conf")

	fmt.Printf("   ... removing '%s'\n", path)

	return os.Remove(path)
}

func (c *Undaemonize) rmdirs(d info) error {
	dir := d.WorkingDirectory

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}

func (c *Undaemonize) firewall(d info) error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Println("   *** WARNING: removing 'uhppoted-rest' to the application firewall")
	fmt.Println("   ***")
	fmt.Println()

	path := d.Executable
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
			fmt.Errorf("ERROR: Failed to remove 'uhppoted-rest' from the application firewall (%v)\n", err)
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
