package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Undaemonize struct {
}

func (c *Undaemonize) Execute(ctx Context) error {
	fmt.Println("   ... undaemonizing")

	path := filepath.Join("/etc/systemd/system", "uhppoted.service")
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

	if err := c.rmdirs(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted unregistered as a systemd service")
	fmt.Println()
	fmt.Println("   Any uhppoted log files can still be found in directory /var/log and should be removed manually")
	fmt.Println()

	return nil
}

func (c *Undaemonize) systemd(path string) error {
	fmt.Printf("   ... stopping uhppoted service\n")
	cmd := exec.Command("systemctl", "stop", "uhppoted")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		fmt.Errorf("ERROR: Failed to stop '%s' (%v)\n", "uhppoted", err)
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
	dir := "/var/uhppoted"

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}

func (c *Undaemonize) Cmd() string {
	return "undaemonize"
}

func (c *Undaemonize) Description() string {
	return "Undaemonizes uhppoted as a service/daemon"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Deregisters uhppoted as a systed service/daemon")
	fmt.Println()
}
