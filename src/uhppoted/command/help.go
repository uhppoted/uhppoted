package uhppoted

import (
	"context"
	"flag"
	"fmt"
	"os"
)

type Help struct {
	service string
	cli     []Command
	run     Command
}

func NewHelp(service string, cli []Command, run Command) *Help {
	return &Help{
		service: service,
		cli:     cli,
		run:     run,
	}
}

func (h *Help) Name() string {
	return "help"
}

func (h *Help) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("help", flag.ExitOnError)
}

func (h *Help) Description() string {
	return "Displays the help information"
}

func (h *Help) Usage() string {
	return ""
}

func (h *Help) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s help <command>\n", h.service)
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("    help          Displays this message")

	for _, c := range h.cli {
		fmt.Printf("    %-13s %s\n", c.FlagSet().Name(), c.Description())
	}
}

func (h *Help) Execute(ctx context.Context) error {
	if len(os.Args) > 2 {
		if os.Args[2] == "commands" {
			h.helpCommands()
			return nil
		}

		if os.Args[2] == h.Name() {
			h.Help()
			return nil
		}

		for _, c := range h.cli {
			if os.Args[2] == c.Name() {
				c.Help()
				return nil
			}
		}

		fmt.Printf("Invalid command: %v. Type 'help commands' to get a list of supported commands\n", flag.Arg(1))
	} else {
		h.usage()
	}

	return nil
}

func (h *Help) usage() {
	fmt.Println()
	fmt.Printf("  Usage: %s <command> [options]\n", h.service)
	fmt.Println()
	fmt.Println("  Defaults to 'run'.")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Printf("    help          Displays this message. For help on a specific command use '%s help <command>'\n", h.service)

	for _, c := range h.cli {
		fmt.Printf("    %-13s %s\n", c.FlagSet().Name(), c.Description())
	}

	fmt.Println()
	fmt.Println(" 'run' options:")

	h.run.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})

	fmt.Println()
}

func (h *Help) helpCommands() {
	fmt.Println()
	fmt.Println("  Supported commands:")

	for _, c := range h.cli {
		fmt.Printf("     %-16s %s\n", c.FlagSet().Name(), c.Usage())
	}

	fmt.Println()
	fmt.Println("     Defaults to 'run'.")
	fmt.Println()
}
