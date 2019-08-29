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
	"uhppoted/eventlog"
	"uhppoted/rest"
)

const (
	LOGFILESIZE = 1
	IDLE        = time.Duration(60 * time.Second)
)

var retries = 0

func main() {
	log.Printf("uhppoted daemon (%d)\n", os.Getpid())

	flag.Parse()

	if err := os.MkdirAll(*dir, os.ModeDir|os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("Error creating working directoryi '%v'", *dir), err)
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

func cleanup(pid string) {
	os.Remove(pid)
}

func run(logfile string, logfilesize int) {
	// ... setup logging

	events := eventlog.Ticker{Filename: logfile, MaxSize: logfilesize}
	logger := log.New(&events, "", 0)

	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)
	rotate := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT)
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

	u := uhppote.UHPPOTE{
		Devices: make(map[uint32]*net.UDPAddr),
		Debug:   true,
	}

	go func() {
		rest.Run(&u)
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
