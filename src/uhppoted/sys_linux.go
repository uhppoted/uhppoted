package main

import (
	"flag"
	"log"
	"os"
)

// LINUX

var configuration = flag.String("config", "/uhppoted/uhppoted.config", "Path for the configuration file")
var dir = flag.String("dir", "/var/uhppoted", "Working directory")
var logfile = flag.String("logfile", "/var/log/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/var/uhppoted/uhppoted.pid", "uhppoted PID file")
var useSyslog = flag.Bool("syslog", false, "Use syslog for event logging")

func sysinit() {
	log.Printf("uhppoted daemon - %s (PID %d)\n", "linux", os.Getpid())
}
