package commands

import (
	"fmt"
)

type Version struct {
	Version string
}

func (c *Version) Execute(ctx Context) error {
	fmt.Printf("%v\n", c.Version)

	return nil
}

func (c *Version) CLI() string {
	return "version"
}

func (c *Version) Description() string {
	return "Displays the current version"
}

func (c *Version) Usage() string {
	return ""
}

func (c *Version) Help() {
	fmt.Println("Displays the uhppoted version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}
