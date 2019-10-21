package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

// WINDOWS (PROVISIONAL - NOT TESTED)

var pwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
var configuration = flag.String("config", filepath.Join(pwd, "uhppoted.config"), "Path for the configuration file")
var dir = flag.String("dir", pwd, "Working directory")
var logfile = flag.String("logfile", filepath.Join(pwd, "logs", "uhppoted.log"), "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", filepath.Join(pwd, "uhppoted.pid"), "uhppoted PID file")

func sysinit() {
	log.Printf("uhppoted daemon - %s (PID %d)\n", "Microsoft Windows", os.Getpid())
}
