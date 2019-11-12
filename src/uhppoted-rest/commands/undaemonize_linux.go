package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	return "Undaemonizes uhppoted-rest as a service/daemon"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest daemonize")
	fmt.Println()
	fmt.Println("    Deregisters uhppoted-rest as a systed service/daemon")
	fmt.Println()
}

func (c *Undaemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... undaemonizing")

	path := filepath.Join("/etc/systemd/system", "uhppoted-rest.service")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... %s does not exist - nothing to do for systemd\n", path)
		return nil
	}

	if err := c.systemd(path); err != nil {
		return err
	}

	if err := c.logrotate(); err != nil {
		return err
	}

	if err := c.rmdirs(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted-rest unregistered as a systemd service")
	fmt.Println()
	fmt.Println("   Log files in directory /var/uhppoted/log and configuration files in /etc/uhppoted should be removed manually")
	fmt.Println()

	return nil
}

func (c *Undaemonize) systemd(path string) error {
	fmt.Printf("   ... stopping uhppoted-rest service\n")
	cmd := exec.Command("systemctl", "stop", "uhppoted-rest")
	out, err := cmd.CombinedOutput()
	if strings.TrimSpace(string(out)) != "" {
		fmt.Printf("   > %s\n", out)
	}
	if err != nil {
		fmt.Errorf("ERROR: Failed to stop '%s' (%v)\n", "uhppoted-rest", err)
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
	path := filepath.Join("/etc/logrotate.d", "uhppoted-rest")

	fmt.Printf("   ... removing '%s'\n", path)

	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Undaemonize) rmdirs() error {
	dir := "/var/uhppoted/rest"

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}
