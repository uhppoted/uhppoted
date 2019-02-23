package commands

import "fmt"

var VERSION = "v0.00.0"

type VersionCommand struct {
	Version string
}

func NewVersionCommand(version string, debug bool) *VersionCommand {
	return &VersionCommand{version}
}

func (c *VersionCommand) Execute() error {
	fmt.Printf("%v\n", c.Version)

	return nil
}
