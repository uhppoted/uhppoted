package commands

import (
	"errors"
	"fmt"
)

type Daemonize struct {
}

func (c *Daemonize) Execute(ctx Context) error {
	return errors.New("uhppoted daemonize: NOT IMPLEMENTED")
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted as a service"
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
