package main

import (
	"flag"
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppote-simulator/commands"
	"github.com/uhppoted/uhppote-simulator/simulator"
	"github.com/uhppoted/uhppote-simulator/simulator/UT0311L04"
	"path/filepath"
)

var options = struct {
	dir   string
	bind  string
	rest  string
	debug bool
}{
	dir:   "./devices",
	bind:  ":60000",
	rest:  ":8000",
	debug: false,
}

func main() {
	flag.StringVar(&options.dir, "devices", "devices", "Specifies the simulation device directory")
	flag.StringVar(&options.bind, "bind", ":60000", "Specifies the bind address for the simulator")
	flag.StringVar(&options.rest, "rest", ":8000", "Specifies the bind address for the REST interface")
	flag.BoolVar(&options.debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	if len(flag.Args()) > 0 && flag.Arg(0) == "help" {
		help()
		return
	}

	if len(flag.Args()) > 0 && flag.Arg(0) == "version" {
		version()
		return
	}

	ctx := simulator.Context{
		Directory:   options.dir,
		BindAddress: options.bind,
		RestAddress: options.rest,
		DeviceList:  simulator.NewDeviceList(load(options.dir)),
	}

	commands.Simulate(&ctx, options.debug)
}

func load(dir string) []simulator.Simulator {
	fmt.Printf("   ... loading devices from '%s'\n", dir)

	devices := map[uint32]simulator.Simulator{}

	list := []struct {
		glob       string
		compressed bool
	}{
		{"*.json", false},
		{"*.json.gz", true},
	}

	for _, g := range list {
		files, err := filepath.Glob(filepath.Join(dir, g.glob))
		if err == nil {
			for _, f := range files {
				s, err := UT0311L04.Load(f, g.compressed)
				if err != nil {
					fmt.Printf("   ... error loading device from file '%s': %v\n", f, err)
				} else {
					deviceID := s.DeviceID()
					if devices[deviceID] == nil {
						devices[deviceID] = s
						fmt.Printf("   ... loaded device  from '%s'\n", f)
					} else {
						fmt.Printf("   ... duplicate serial number %v in device file '%s' - using device loaded from '%s'\n", deviceID, f, devices[deviceID].FilePath())
					}
				}
			}
		}
	}

	l := make([]simulator.Simulator, 0)

	for _, s := range devices {
		l = append(l, s)
	}

	return l
}

func version() {
	fmt.Printf("%v\n", uhppote.VERSION)
}

func help() {
	fmt.Println()
	fmt.Println("  Usage: uhppote-simulator [options] <command>")
	fmt.Println()
	fmt.Println("  By default, the application will run in 'simulation' mode if a command is not specified.")
	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println()
	fmt.Println("    version          Displays the simulator version")
	fmt.Println("    help             Displays this message")
	fmt.Println("                     For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --devices <directory>  Sets the directory to which to load and save simulator device files")
	fmt.Println("    --debug                Displays vaguely useful internal information")
	fmt.Println()
}
