package commands

import (
	"flag"
	"fmt"
)

type Version struct {
	version string
	flagset *flag.FlagSet
}

var version = Version{
	version: VERSION,
	flagset: flag.NewFlagSet("version", flag.ExitOnError),
}

func (c *Version) FlagSet() *flag.FlagSet {
	return c.flagset
}

func (c *Version) Parse(args []string) error {
	flagset := c.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'version' command implementation without a flagset: %#v", c))
	}

	return flagset.Parse(args)
}

func (c *Version) Execute(ctx Context) error {
	fmt.Printf("%v\n", c.version)

	return nil
}

func (c *Version) Cmd() string {
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
