package commands

import (
	"context"
	"flag"
	"fmt"
)

var DUMP = Dump{
	config: "/etc/uhppoted/uhppoted.conf",
}

type Dump struct {
	config string
}

func (d *Dump) Name() string {
	return "config"
}

func (d *Dump) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("config", flag.ExitOnError)
}

func (d *Dump) Description() string {
	return fmt.Sprintf("Displays all the configuration information for %s", SERVICE)
}

func (d *Dump) Usage() string {
	return ""
}

func (d *Dump) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s config\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Displays all the configuration information for %s\n", SERVICE)
	fmt.Println()
}

func (d *Dump) Execute(ctx context.Context) error {
	if err := dump(d.config); err != nil {
		return err
	}

	return nil
}
