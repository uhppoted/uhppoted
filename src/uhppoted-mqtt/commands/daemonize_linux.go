package commands

import (
	"context"
	"flag"
	"fmt"
)

var DAEMONIZE = Daemonize{}

type Daemonize struct {
}

func NewDaemonize() *Daemonize {
	return &Daemonize{}
}

func (d *Daemonize) Name() string {
	return "daemonize"
}

func (d *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (d *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (d *Daemonize) Usage() string {
	return ""
}

func (d *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Daemonizes %s as a service/daemon that runs on startup\n", SERVICE)
	fmt.Println()
}

func (d *Daemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... daemonizing")

	return nil
}
