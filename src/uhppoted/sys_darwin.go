package main

import (
	"flag"
	"log"
	"os"
)

// DARWIN

var dir = flag.String("dir", "/usr/local/var/uhppoted", "Working directory")
var logfile = flag.String("logfile", "/usr/local/var/uhppoted/logs/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/usr/local/var/uhppoted/uhppoted.pid", "uhppoted PID file")
var useSyslog = flag.Bool("syslog", false, "Use syslog for event logging")

func sysinit() {
	log.Printf("uhppoted daemon %s - %s (PID %d)\n", VERSION, "MacOS", os.Getpid())
}