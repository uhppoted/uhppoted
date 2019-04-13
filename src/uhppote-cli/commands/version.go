package commands

import (
	"fmt"
	"uhppote"
)

var VERSION = "v0.00.0"

type VersionCommand struct {
	Version string
}

func (c *VersionCommand) Execute(u *uhppote.UHPPOTE) error {
	fmt.Printf("%v\n", c.Version)

	return nil
}

func (c *VersionCommand) CLI() string {
	return "version"
}

func (c *VersionCommand) Help() {
	fmt.Println("Displays the uhppote-cli version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}
