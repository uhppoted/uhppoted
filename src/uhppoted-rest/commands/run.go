package commands

import (
	"context"
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
	"uhppoted-rest/config"
	"uhppoted-rest/rest"
)

type status struct {
	touched time.Time
	status  types.Status
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
}

type state struct {
	started time.Time

	healthcheck struct {
		touched *time.Time
		alerted bool
	}

	devices struct {
		status sync.Map
		errors sync.Map
	}
}

const (
	SERVICE = `uhppoted-rest`
	IDLE    = time.Duration(60 * time.Second)
	IGNORE  = time.Duration(5 * time.Minute)
	DELTA   = 60
	DELAY   = 30
)

func (c *Run) Name() string {
	return "run"
}

func (c *Run) Description() string {
	return "Runs the uhppoted-rest daemon/service until terminated by the system service manager"
}

func (c *Run) Usage() string {
	return "uhppoted-rest [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]"
}

func (c *Run) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest <options>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	c.FlagSet().VisitAll(func(f *flag.Flag) {
		fmt.Printf("    --%-12s %s\n", f.Name, f.Usage)
	})
	fmt.Println()
}

func (r *Run) execute(ctx context.Context, f func(*config.Config) error) error {
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
	logger.Printf("START")

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

	logger.Printf("STOP")
}

func (r *Run) listen(c *config.Config, logger *log.Logger, interrupt chan os.Signal) error {
	s := state{
		started: time.Now(),

		healthcheck: struct {
			touched *time.Time
			alerted bool
		}{
			touched: nil,
			alerted: false,
		},
		devices: struct {
			status sync.Map
			errors sync.Map
		}{
			status: sync.Map{},
			errors: sync.Map{},
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
			if err := watchdog(&u, &s, logger); err != nil {
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
}

func healthcheck(u *uhppote.UHPPOTE, st *state, l *log.Logger) {
	l.Printf("health-check")

	now := time.Now()
	devices := make(map[uint32]bool)

	found, err := u.FindDevices()
	if err != nil {
		l.Printf("WARN  'keep-alive' error: %v", err)
	}

	if found != nil {
		for _, id := range found {
			devices[uint32(id.SerialNumber)] = true
		}
	}

	for id, _ := range u.Devices {
		devices[id] = true
	}

	for id, _ := range devices {
		s, err := u.GetStatus(id)
		if err == nil {
			st.devices.status.Store(id, status{
				touched: now,
				status:  *s,
			})
		}
	}

	st.healthcheck.touched = &now
}

func watchdog(u *uhppote.UHPPOTE, st *state, l *log.Logger) error {
	warnings := 0
	errors := 0
	healthCheckRunning := false
	now := time.Now()
	seconds, _ := time.ParseDuration("1s")

	// Verify health-check

	dt := time.Since(st.started).Round(seconds)
	if st.healthcheck.touched != nil {
		dt = time.Since(*st.healthcheck.touched)
		if int64(math.Abs(dt.Seconds())) < DELAY {
			healthCheckRunning = true
		}
	}

	if int64(math.Abs(dt.Seconds())) > DELAY {
		errors += 1
		if !st.healthcheck.alerted {
			l.Printf("ERROR 'health-check' subsystem has not run since %s (%s)", types.DateTime(st.started), dt)
			st.healthcheck.alerted = true
		}
	} else {
		if st.healthcheck.alerted {
			l.Printf("INFO  'health-check' subsystem is running")
			st.healthcheck.alerted = false
		}
	}

	// Verify configured devices

	if healthCheckRunning {
		for id, _ := range u.Devices {
			alerted := alerts{
				missing:      false,
				unexpected:   false,
				touched:      false,
				synchronized: false,
			}

			if v, found := st.devices.errors.Load(id); found {
				alerted.missing = v.(alerts).missing
				alerted.unexpected = v.(alerts).unexpected
				alerted.touched = v.(alerts).touched
				alerted.synchronized = v.(alerts).synchronized
			}

			if _, found := st.devices.status.Load(id); !found {
				errors += 1
				if !alerted.missing {
					l.Printf("ERROR UTC0311-L0x %s device not found", types.SerialNumber(id))
					alerted.missing = true
				}
			}

			if v, found := st.devices.status.Load(id); found {
				touched := v.(status).touched
				t := time.Time(v.(status).status.SystemDateTime)
				dt := time.Since(t).Round(seconds)
				dtt := int64(math.Abs(time.Since(touched).Seconds()))

				if alerted.missing {
					l.Printf("ERROR UTC0311-L0x %s present", types.SerialNumber(id))
					alerted.missing = false
				}

				if now.After(touched.Add(IDLE)) {
					errors += 1
					if !alerted.touched {
						l.Printf("ERROR UTC0311-L0x %s no response for %s", types.SerialNumber(id), time.Since(touched).Round(seconds))
						alerted.touched = true
						alerted.synchronized = false
					}
				} else {
					if alerted.touched {
						l.Printf("INFO  UTC0311-L0x %s connected", types.SerialNumber(id))
						alerted.touched = false
					}
				}

				if dtt < DELTA/2 {
					if int64(math.Abs(dt.Seconds())) > DELTA {
						errors += 1
						if !alerted.synchronized {
							l.Printf("ERROR UTC0311-L0x %s system time not synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
							alerted.synchronized = true
						}
					} else {
						if alerted.synchronized {
							l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
							alerted.synchronized = false
						}
					}
				}
			}

			st.devices.errors.Store(id, alerted)
		}
	}

	// Any unexpected devices?

	st.devices.status.Range(func(key, value interface{}) bool {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
		}

		if v, found := st.devices.errors.Load(key); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		for id, _ := range u.Devices {
			if id == key {
				if alerted.unexpected {
					l.Printf("ERROR UTC0311-L0x %s added to configuration", types.SerialNumber(key.(uint32)))
					alerted.unexpected = false
					st.devices.errors.Store(id, alerted)
				}

				return true
			}
		}

		touched := value.(status).touched
		t := time.Time(value.(status).status.SystemDateTime)
		dt := time.Since(t).Round(seconds)
		dtt := int64(math.Abs(time.Since(touched).Seconds()))

		if now.After(touched.Add(IGNORE)) {
			st.devices.status.Delete(key)
			st.devices.errors.Delete(key)

			if alerted.unexpected {
				l.Printf("WARN  UTC0311-L0x %s disappeared", types.SerialNumber(key.(uint32)))
			}
		} else {
			warnings += 1
			if !alerted.unexpected {
				l.Printf("WARN  UTC0311-L0x %s unexpected device", types.SerialNumber(key.(uint32)))
				alerted.unexpected = true
			}

			if now.After(touched.Add(IDLE)) {
				warnings += 1
				if !alerted.touched {
					l.Printf("WARN  UTC0311-L0x %s no response for %s", types.SerialNumber(key.(uint32)), time.Since(touched).Round(seconds))
					alerted.touched = true
					alerted.synchronized = false
				}
			} else {
				if alerted.touched {
					l.Printf("INFO  UTC0311-L0x %s connected", types.SerialNumber(key.(uint32)))
					alerted.touched = false
				}
			}

			if dtt < DELTA/2 {
				if int64(math.Abs(dt.Seconds())) > DELTA {
					warnings += 1
					if !alerted.synchronized {
						l.Printf("WARN  UTC0311-L0x %s system time not synchronized: %s (%s)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
						alerted.synchronized = true
					}
				} else {
					if alerted.synchronized {
						l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
						alerted.synchronized = false
					}
				}
			}

			st.devices.errors.Store(key, alerted)
		}

		return true
	})

	// 'k, done

	if errors > 0 {
		l.Printf("watchdog: ERROR")
	} else if warnings > 0 {
		l.Printf("watchdog: WARN")
	} else {
		l.Printf("watchdog: OK")
	}

	return nil
}
