package commands

import (
	"fmt"
	"uhppote"
)

var VERSION = "v0.00.0"

type VersionCommand struct {
	Version string
}

func NewVersionCommand(version string) *VersionCommand {
	return &VersionCommand{version}
}

func (c *VersionCommand) Execute(u *uhppote.UHPPOTE) error {
	fmt.Printf("%v\n", c.Version)

	return nil
}
