package commands

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Undaemonize struct {
}

func (c *Undaemonize) Execute(ctx Context) error {
	return errors.New("uhppoted undaemonize: NOT IMPLEMENTED")
}

func (c *Undaemonize) CLI() string {
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
