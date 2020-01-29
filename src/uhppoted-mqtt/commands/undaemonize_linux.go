package commands

import (
	"context"
	"flag"
	"fmt"
)

var UNDAEMONIZE = Undaemonize{}

type Undaemonize struct {
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{}
}

func (u *Undaemonize) Name() string {
	return "undaemonize"
}

func (u *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (u *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters %s as a service/daemon", SERVICE)
}

func (u *Undaemonize) Usage() string {
	return ""
}

func (u *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s from launchd as a service/daemon", SERVICE)
	fmt.Println()
}

func (u *Undaemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... undaemonizing")

	return nil
}
