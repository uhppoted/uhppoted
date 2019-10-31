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

var configuration = flag.String("config", "/etc/uhppoted/uhppoted.conf", "Path for the configuration file")
var dir = flag.String("dir", "/var/uhppoted", "Working directory")
var logfile = flag.String("logfile", "/var/log/uhppoted/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/var/uhppoted/uhppoted.pid", "uhppoted PID file")

func (c *Run) Execute(ctx Context) error {
	log.Printf("uhppoted daemon - %s (PID %d)\n", "Linux", os.Getpid())

	return execute(ctx)
}

func start(c *config.Config, logfile string, logfilesize int) {
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

	run(c, logger)
}
