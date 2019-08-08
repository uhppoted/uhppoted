package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"uhppote-simulator/commands"
	"uhppote-simulator/entities"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

var cli = []commands.Command{
	&commands.VersionCommand{VERSION},
}

var VERSION = "v0.03.0"

var options = struct {
	dir   string
	debug bool
}{
	dir:   "./devices",
	debug: false,
}

func main() {
	flag.StringVar(&options.dir, "devices", "devices", "Specifies the simulation device directory")
	flag.BoolVar(&options.debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd != nil {
		cmd.Execute(options.dir)
		return
	}

	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		help()
		return
	}

	ctx := simulator.Context{
		Directory: options.dir,
		DeviceList: simulator.DeviceList{
			TxQ:        make(chan entities.Message, 8),
			Simulators: load(options.dir),
		},
	}

	simulate(&ctx)
}

func parse() (commands.Command, error) {
	var cmd commands.Command = nil
	var err error = nil

	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.CLI() == flag.Arg(0) {
				cmd = c
			}
		}
	}

	return cmd, err
}

func load(dir string) []*simulator.Simulator {
	fmt.Printf("   ... loading devices from '%s'\n", dir)

	devices := map[types.SerialNumber]*simulator.Simulator{}

	list := []struct {
		glob       string
		compressed bool
	}{
		{"*.json", false},
		{"*.json.gz", true},
	}

	for _, g := range list {
		files, err := filepath.Glob(path.Join(dir, g.glob))
		if err == nil {
			for _, f := range files {
				s, err := simulator.Load(f, g.compressed)
				if err != nil {
					fmt.Printf("   ... error loading device from file '%s': %v\n", f, err)
				} else {
					if devices[s.SerialNumber] == nil {
						devices[s.SerialNumber] = s
						fmt.Printf("   ... loaded device  from '%s'\n", f)
					} else {
						fmt.Printf("   ... duplicate serial number %v in device file '%s' - using device loaded from '%s'\n", s.SerialNumber, f, devices[s.SerialNumber].File)
					}
				}
			}
		}
	}

	fmt.Println()

	l := make([]*simulator.Simulator, 0)

	for _, s := range devices {
		l = append(l, s)
	}

	return l
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

func usage() {
	fmt.Println()
	fmt.Println("  Usage: uhppote-simulator [options] <command>")
	fmt.Println()
	fmt.Println("  By default, the application will run in 'simulation' mode if a command is not specified.")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    help             Displays this message")
	fmt.Println("                     For help on a specific command use 'uhppote-cli help <command>'")

	for _, c := range cli {
		fmt.Printf("    %-16s %s\n", c.CLI(), c.Description())
	}

	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --devices   Sets the directory to which to load and save simulator device files")
	fmt.Println("    --debug     Displays vaguely useful internal information")
	fmt.Println()
}

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()

	for _, c := range cli {
		fmt.Printf(" %-16s %s\n", c.CLI(), c.Usage())
	}

	fmt.Println()
}
