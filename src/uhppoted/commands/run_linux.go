package commands

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"uhppoted/config"
	"uhppoted/eventlog"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	debug         bool
}

var runCmd = Run{
	configuration: "/etc/uhppoted/uhppoted.conf",
	dir:           "/var/uhppoted",
	pidFile:       "/var/uhppoted/uhppoted.pid",
	logFile:       "/var/log/uhppoted/uhppoted.log",
	logFileSize:   10,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Path for the configuration file")
	flagset.StringVar(&r.dir, "dir", r.dir, "Working directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "uhppoted PID file")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "uhppoted log file")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "uhppoted log file size")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays vaguely useful internal information")

	return flagset
}

func (r *Run) Execute(ctx Context) error {
	log.Printf("uhppoted daemon - %s (PID %d)\n", "Linux", os.Getpid())

	return r.execute(ctx)
}

func (r *Run) start(c *config.Config, logfile string, logfilesize int) {
	// ... setup logging

	events := eventlog.Ticker{Filename: logfile, MaxSize: logfilesize}
	logger := log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
	rotate := make(chan os.Signal, 1)

	signal.Notify(rotate, syscall.SIGHUP)

	go func() {
		for {
			<-rotate
			log.Printf("Rotating uhppoted log file '%s'\n", logfile)
			events.Rotate()
		}
	}()

	r.run(c, logger)
}
