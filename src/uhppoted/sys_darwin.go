package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"uhppoted/config"
	"uhppoted/eventlog"
)

// DARWIN

var configuration = flag.String("config", "/usr/local/etc/com.github.twystd.uhppoted/uhppoted.conf", "Path for the configuration file")
var dir = flag.String("dir", "/usr/local/var/com.github.twystd.uhppoted", "Working directory")
var logfile = flag.String("logfile", "/usr/local/var/com.github.twystd.uhppoted/logs/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/usr/local/var/com.github.twystd.uhppoted/uhppoted.pid", "uhppoted PID file")

func sysinit() {
	log.Printf("uhppoted daemon %s - %s (PID %d)\n", VERSION, "MacOS", os.Getpid())
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
