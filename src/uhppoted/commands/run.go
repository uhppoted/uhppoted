package commands

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uhppote"
	"uhppote/types"
	"uhppoted/config"
	"uhppoted/rest"
)

type state struct {
	devices struct {
		touched sync.Map
		status  sync.Map
		errors  sync.Map
	}
}

type alerts struct {
	touched      bool
	synchronized bool
}

const (
	IDLE  = time.Duration(60 * time.Second)
	DELTA = 60
)

func (c *Run) Parse(args []string) error {
	flagset := c.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'run' command implementation without a flagset: %#v", c))
	}

	return flagset.Parse(args)
}

func (c *Run) Description() string {
	return "Runs the uhppoted daemon/service until terminated by the system service manager"
}

func (c *Run) Usage() string {
	return "uhppoted [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
}

func (c *Run) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted <options>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	runCmd.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})
	fmt.Println()
}

func (r *Run) execute(ctx Context, f func(*config.Config) error) error {
	conf := config.NewConfig()
	if err := conf.Load(r.configuration); err != nil {
		log.Printf("\n   WARN:  Could not load configuration (%v)\n\n", err)
	}

	if err := os.MkdirAll(r.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", r.dir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	if err := ioutil.WriteFile(r.pidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create pid file: %v\n", err)
	}

	defer func() {
		os.Remove(r.pidFile)
	}()

	return f(conf)
}

func (r *Run) run(c *config.Config, logger *log.Logger) {
	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ... listen forever

	for {
		err := r.listen(c, logger, interrupt)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		log.Printf("exit\n")
		break
	}
}

func (r *Run) listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {

	s := state{
		devices: struct {
			touched sync.Map
			status  sync.Map
			errors  sync.Map
		}{
			touched: sync.Map{},
			status:  sync.Map{},
			errors:  sync.Map{},
		},
	}

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            r.debug,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
	}

	// ... REST task

	restd := rest.RestD{
		HttpEnabled:        c.REST.HttpEnabled,
		HttpPort:           c.REST.HttpPort,
		HttpsEnabled:       c.REST.HttpsEnabled,
		HttpsPort:          c.REST.HttpsPort,
		TLSKeyFile:         c.REST.TLSKeyFile,
		TLSCertificateFile: c.REST.TLSCertificateFile,
		CACertificateFile:  c.REST.CACertificateFile,
		CORSEnabled:        c.REST.CORSEnabled,
		OpenApi: rest.OpenApi{
			Enabled:   c.OpenApi.Enabled,
			Directory: c.OpenApi.Directory,
		},
	}

	go func() {
		restd.Run(&u, logger)
	}()

	defer rest.Close()

	// ... health-check task

	k := time.NewTicker(15 * time.Second)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck(&u, &s, logger)
		}
	}()

	// ... wait until interrupted/closed

	closed := make(chan struct{})
	w := time.NewTicker(5 * time.Second)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog(&s, logger); err != nil {
				return err
			}

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

func healthcheck(u *uhppote.UHPPOTE, s *state, l *log.Logger) {
	l.Printf("health-check")
	now := time.Now()
	devices, err := u.FindDevices()

	if err != nil {
		l.Printf("WARN  'keep-alive' error: %v", err)
		return
	}

	for _, device := range devices {
		s.devices.touched.Store(device.SerialNumber, now)
	}

	for _, device := range devices {
		status, err := u.GetStatus(uint32(device.SerialNumber))
		if err == nil {
			s.devices.status.Store(device.SerialNumber, *status)
		}
	}
}

func watchdog(s *state, l *log.Logger) error {
	ok := true
	now := time.Now()
	seconds, _ := time.ParseDuration("1s")
	alerted := alerts{
		touched:      false,
		synchronized: false,
	}

	s.devices.touched.Range(func(key, value interface{}) bool {
		touched := value.(time.Time)
		timeout := touched.Add(IDLE)

		if v, found := s.devices.errors.Load(key); found {
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		if now.After(timeout) {
			ok = false
			if !alerted.touched {
				l.Printf("ERROR UTC0311-L0x %s no response for %s", key, time.Since(touched).Round(seconds))
				alerted.touched = true
			}
		} else {
			if alerted.touched {
				l.Printf("INFO  UTC0311-L0x %s reconnected", key)
				alerted.touched = false
			}
		}

		s.devices.errors.Store(key, alerted)

		return true
	})

	s.devices.status.Range(func(key, value interface{}) bool {
		status := value.(types.Status)
		t := time.Time(status.SystemDateTime)
		dt := time.Since(t).Round(seconds)

		if v, found := s.devices.errors.Load(key); found {
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		if int64(math.Abs(dt.Seconds())) > DELTA {
			ok = false
			if !alerted.synchronized {
				l.Printf("ERROR UTC0311-L0x %s system time not synchronized: %s (%s)", key, status.SystemDateTime, dt)
				alerted.synchronized = true
			}
		} else {
			if alerted.synchronized {
				l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", key, status.SystemDateTime, dt)
				alerted.synchronized = false
			}
		}

		s.devices.errors.Store(key, alerted)

		return true
	})

	if ok {
		l.Printf("watchdog: OK")
	} else {
		l.Printf("watchdog: ERROR")
	}

	return nil
}
