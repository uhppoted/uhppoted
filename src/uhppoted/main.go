package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uhppote"
	"uhppoted/commands"
	"uhppoted/eventlog"
	"uhppoted/rest"
)

const (
	LOGFILESIZE = 1
	IDLE        = time.Duration(60 * time.Second)
)

var cli = []commands.Command{
	&commands.Install{},
}

var VERSION = "v0.04.0"
var retries = 0

func main() {
	flag.Parse()

	cmd, err := parse()
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	ctx := commands.Context{}

	if cmd != nil {
		if err = cmd.Execute(ctx); err != nil {
			fmt.Printf("\nERROR: %v\n\n", err)
			os.Exit(1)
		}

		return
	}

	// ... default to 'run'

	sysinit()

	if err := os.MkdirAll(*dir, os.ModeDir|os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("Error creating working directory '%v'", *dir), err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
		log.Fatal("Error creating pid file: %v\n", err)
	}

	defer cleanup(*pidFile)

	// ... use syslog for console logging?

	if *useSyslog {
		logger, err := syslog.New(syslog.LOG_NOTICE, "uhppoted")

		if err != nil {
			log.Fatal("Error opening syslog: ", err)
			return
		}

		log.SetOutput(logger)
	}

	run(*logfile, *logfilesize)
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

func cleanup(pid string) {
	os.Remove(pid)
}

func run(logfile string, logfilesize int) {
	// ... setup logging

	events := eventlog.Ticker{Filename: logfile, MaxSize: logfilesize}
	logger := log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)

	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)
	rotate := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(rotate, syscall.SIGHUP)

	go func() {
		for {
			<-rotate
			log.Printf("Rotating uhppoted log file '%s'\n", logfile)
			events.Rotate()
		}
	}()

	// ... listen forever

	for {
		err := listen(logger, interrupt)

		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}
}

func listen(logger *log.Logger, interrupt chan os.Signal) error {
	// ... listen

	log.Printf("... listening")

	address, _ := net.ResolveUDPAddr("udp", "0.0.0.0:60001")

	u := uhppote.UHPPOTE{
		BindAddress: address,
		Devices:     make(map[uint32]*net.UDPAddr),
		Debug:       true,
	}

	go func() {
		rest.Run(&u, logger)
	}()

	defer rest.Close()

	touched := time.Now()
	closed := make(chan struct{})

	// ... wait until interrupted/closed

	k := time.NewTicker(15 * time.Second)
	tick := time.NewTicker(5 * time.Second)

	defer k.Stop()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if err := watchdog(touched); err != nil {
				return err
			}

		case <-k.C:
			log.Printf("... keep-alive")
			keepalive()

		case <-interrupt:
			log.Printf("... interrupt")
			return nil

		case <-closed:
			log.Printf("... closed")
			return errors.New("Server error")
		}
	}

	log.Printf("... exit")
	return nil
}
