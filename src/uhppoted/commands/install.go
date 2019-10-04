package commands

import (
	"fmt"
)

type Install struct {
}

func (c *Install) Execute(ctx Context) error {
	fmt.Println("...... installing")
	return nil
}

func (c *Install) CLI() string {
	return "install"
}

func (c *Install) Description() string {
	return "Installs uhppoted as a service/daemon"
}

func (c *Install) Usage() string {
	return ""
}

func (c *Install) Help() {
	fmt.Println("Usage: uhppoted install")
	fmt.Println()
	fmt.Println(" Installs uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}
