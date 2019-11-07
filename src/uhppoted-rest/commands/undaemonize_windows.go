package commands

import (
	"flag"
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
		name:        "uhppoted-rest",
		description: "uhppoted-rest Service Interface to UTO311-L0x devices",
	}
}

func (c *Undaemonize) Name() string {
	return "undaemonize"
}

func (c *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (c *Undaemonize) Execute(ctx Context) error {
	fmt.Println("   ... undaemonizing")

	dir := workdir()
	if err := c.unregister(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted-rest unregistered as a Windows system service")
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

func (c *Undaemonize) Description() string {
	return "Deregisters the uhppoted-rest service"
}

func (c *Undaemonize) Usage() string {
	return ""
}

func (c *Undaemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest undaemonize")
	fmt.Println()
	fmt.Println("    Deregisters uhppoted-rest as a Windows service")
	fmt.Println()
}
