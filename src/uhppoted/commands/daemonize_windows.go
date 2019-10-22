package commands

import (
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"os"
)

type Daemonize struct {
	name        string
	description string
}

func NewDaemonize() *Daemonize {
	return &Daemonize{
		name:        "uhppoted",
		description: "uhppoted Service Interface to UTO311-L0x devices",
	}
}

func (c *Daemonize) Parse(args []string) error {
	return nil
}

func (c *Daemonize) Execute(ctx Context) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(c.name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", c.name)
	}

	s, err = m.CreateService(c.name, executable, mgr.Config{DisplayName: c.description}, "is", "auto-started")
	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(c.name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}

	return nil
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted as a Windows service"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Registers uhppoted as a windows Service that runs on startup")
	fmt.Println()
}
