package commands

import (
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

type Undaemonize struct {
	name        string
	description string
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{
		name:        "uhppoted",
		description: "uhppoted Service Interface to UTO311-L0x devices",
	}
}

func (c *Undaemonize) Parse(args []string) error {
	return nil
}

func (c *Undaemonize) Execute(ctx Context) error {
	fmt.Println("   ... undaemonizing")

	dir := workdir()
	if err := c.unregister(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted unregistered as a Windows system service")
	fmt.Println()
	fmt.Printf("   Log files and configuration files in directory %s should be removed manually", dir)
	fmt.Println()

	return nil
}

func (c *Undaemonize) unregister() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(c.name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", c.name)
	}

	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}

	err = eventlog.Remove(c.name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}

	return nil
}

func (c *Undaemonize) Cmd() string {
	return "undaemonize"
}

func (c *Undaemonize) Description() string {
	return "Deregisters the uhppoted service"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println("Usage: uhppoted undaemonize")
	fmt.Println()
	fmt.Println(" Deregisters uhppoted as a Windows service")
	fmt.Println()
}
