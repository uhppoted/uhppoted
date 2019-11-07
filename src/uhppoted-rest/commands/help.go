package commands

import (
	"flag"
	"fmt"
	"os"
)

type Help struct {
}

func (c *Help) Name() string {
	return "help"
}

func (c *Help) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("help", flag.ExitOnError)
}

func (c *Help) Description() string {
	return "Displays the current version"
}

func (c *Help) Usage() string {
	return ""
}

func (c *Help) Help() {
	fmt.Println()
	fmt.Println("Displays the uhppoted-rest version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}

func (c *Help) Execute(ctx Context) error {
	if len(os.Args) > 2 {
		if os.Args[2] == "commands" {
			helpCommands()
			return nil
		}

		for _, c := range cli {
			flagset := c.FlagSet()
			if flagset.Name() == os.Args[2] {
				c.Help()
				return nil
			}
		}

		fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
	} else {
		usage()
	}

	return nil
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest <command> [options]")
	fmt.Println()
	fmt.Println("  Defaults to 'run'.")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help          Displays this message. For help on a specific command use 'uhppoted-rest help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-13s %s\n", c.FlagSet().Name(), c.Description())
	}

	fmt.Println()
	fmt.Println(" 'run' options:")

	runCmd.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})

	fmt.Println()
}

func helpCommands() {
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range cli {
		fmt.Printf("     %-16s %s\n", c.FlagSet().Name(), c.Usage())
	}

	fmt.Println()
	fmt.Println("     Defaults to 'run'.")
	fmt.Println()
}
