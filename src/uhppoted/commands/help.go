package commands

import (
	"flag"
	"fmt"
)

type Help struct {
	commands []Command
}

func (c *Help) Execute(ctx Context) error {
	if len(flag.Args()) > 1 {
		if flag.Arg(1) == "commands" {
			helpCommands()
			return nil
		}

		for _, c := range cli {
			if c.CLI() == flag.Arg(1) {
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

func (c *Help) CLI() string {
	return "help"
}

func (c *Help) Description() string {
	return "Displays the current version"
}

func (c *Help) Usage() string {
	return ""
}

func (c *Help) Help() {
	fmt.Println("Displays the uhppoted version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted [options] <command>")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help          Displays this message. For help on a specific command use 'uhppoted help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-13s %s\n", c.CLI(), c.Description())
	}

	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --dir         Sets the working directory")
	fmt.Println("    --logfile     Sets the log file path")
	fmt.Println("    --logfilesize Sets the log file size before forcing a log rotate")
	fmt.Println("    --pid         Sets the PID file path")
	fmt.Println("    --syslog      Writes log information to the syslog")
	fmt.Println("    --debug       Displays vaguely useful internal information")
	fmt.Println()
}

func help() {
	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		if len(flag.Args()) > 1 {

			if flag.Arg(1) == "commands" {
				helpCommands()
				return
			}

			for _, c := range cli {
				if c.CLI() == flag.Arg(1) {
					c.Help()
					return
				}
			}

			fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
			return
		}
	}

	usage()
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(" %-16s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}
