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

type Daemonize struct {
}

func (c *Daemonize) Execute(ctx Context) error {
	return errors.New("uhppoted install: NOT IMPLEMENTED")
}

func (c *Daemonize) CLI() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Daemonizes uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Daemonizes uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}
