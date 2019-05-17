package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"uhppote-simulator/commands"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

var cli = []commands.Command{
	&commands.VersionCommand{VERSION},
	&commands.NewDeviceCommand{},
}

var VERSION = "v0.00.0"

var options = struct {
	dir   string
	debug bool
}{
	dir:   "./devices",
	debug: false,
}

var simulators = []*simulator.Simulator{}

func main() {
	flag.StringVar(&options.dir, "devices", "./devices", "Specifies the simulation device directory")
	flag.BoolVar(&options.debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\n   ERROR: %v\n\n", err)
		os.Exit(1)
	}

	if cmd != nil {
		cmd.Execute()
		return
	}

	simulate()
}

func parse() (commands.Command, error) {
	var cmd commands.Command = nil
	var err error = nil

	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		default:
			for _, c := range cli {
				if c.CLI() == flag.Arg(0) {
					cmd = c
				}
			}
		}
	}

	return cmd, err
}

func load(dir string) []*simulator.Simulator {
	if options.debug {
		fmt.Printf("   ... loading devices from '%s'\n", dir)
	}

	devices := map[types.SerialNumber]*simulator.Simulator{}

	files, err := filepath.Glob(path.Join(dir, "*.gz"))
	if err == nil {
		for _, f := range files {
			s, err := simulator.LoadGZ(f)
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

	files, err = filepath.Glob(path.Join(dir, "*.json"))
	if err == nil {
		for _, f := range files {
			s, err := simulator.Load(f)
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

	fmt.Println()

	l := make([]*simulator.Simulator, len(simulators))

	for _, s := range devices {
		l = append(l, s)
	}

	return l
}
