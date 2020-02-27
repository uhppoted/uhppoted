package commands

import (
	"context"
	"flag"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"github.com/uhppoted/uhppoted/src/uhppoted-rest/config"
	"github.com/uhppoted/uhppoted/src/uhppoted-rest/eventlog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	console       bool
	debug         bool
}

var RUN = Run{
	configuration: "/etc/uhppoted/uhppoted.conf",
	dir:           "/var/uhppoted",
	pidFile:       "/var/uhppoted/uhppoted-rest.pid",
	logFile:       "/var/log/uhppoted/uhppoted-rest.log",
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Sets the configuration file path")
	flagset.StringVar(&r.dir, "dir", r.dir, "Work directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "Sets the service PID file path")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "Sets the log file path")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&r.console, "console", r.console, "Writes log entries to stdout")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays vaguely useful internal information")

	return flagset
}

func (r *Run) Execute(ctx context.Context) error {
	log.Printf("uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "Linux", os.Getpid())

	f := func(c *config.Config) error {
		return r.exec(c)
	}

	return r.execute(ctx, f)
}

func (r *Run) exec(c *config.Config) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !r.console {
		events := eventlog.Ticker{Filename: r.logFile, MaxSize: r.logFileSize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Printf("Rotating uhppoted-rest log file '%s'\n", r.logFile)
				events.Rotate()
			}
		}()
	}

	r.run(c, logger)

	return nil
}
