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
	return "daemonize"
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
