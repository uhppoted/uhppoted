package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uhppote"
	"uhppoted/commands"
	"uhppoted/config"
	"uhppoted/rest"
)

const (
	LOGFILESIZE = 1
	IDLE        = time.Duration(60 * time.Second)
)

var VERSION = "v0.04.0"
var retries = 0

func main() {
	flag.Parse()

	cmd, err := commands.Parse()
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if cmd != nil {
		ctx := commands.Context{}
		if err = cmd.Execute(ctx); err != nil {
			fmt.Printf("\nERROR: %v\n\n", err)
			os.Exit(1)
		}

		return
	}

	// ... default to 'run'

	sysinit()

	conf := config.NewConfig()
	if err := conf.Load(*configuration); err != nil {
		fmt.Printf("\n   WARN:  Could not load configuration (%v)\n\n", err)
	}

	if err := os.MkdirAll(*dir, os.ModeDir|os.ModePerm); err != nil {
		log.Fatal(fmt.Sprintf("ERROR: unable to create working directory '%v'", *dir), err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(*pidFile, []byte(pid), 0644); err != nil {
		log.Fatal(fmt.Sprintf("ERROR: unable to create pid file: %v\n", err))
	}

	defer cleanup(*pidFile)

	start(conf, *logfile, *logfilesize)
}

func cleanup(pid string) {
	os.Remove(pid)
}

func run(c *config.Config, logger *log.Logger) {
	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ... listen forever

	for {
		err := listen(c, logger, interrupt)

		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}
}

func listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {
	// ... listen

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            true,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
	}

	restd := rest.RestD{
		HttpEnabled:        c.REST.HttpEnabled,
		HttpPort:           c.REST.HttpPort,
		HttpsEnabled:       c.REST.HttpsEnabled,
		HttpsPort:          c.REST.HttpsPort,
		TLSKeyFile:         c.REST.TLSKeyFile,
		TLSCertificateFile: c.REST.TLSCertificateFile,
		CACertificateFile:  c.REST.CACertificateFile,
		CORSEnabled:        c.REST.CORSEnabled,
	}

	go func() {
		restd.Run(&u, logger)
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
			logger.Printf("... keep-alive")
			keepalive()

		case <-interrupt:
			logger.Printf("... interrupt")
			return nil

		case <-closed:
			logger.Printf("... closed")
			return errors.New("Server error")
		}
	}

	logger.Printf("... exit")
	return nil
}
