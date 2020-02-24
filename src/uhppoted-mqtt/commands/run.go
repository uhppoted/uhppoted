package commands

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uhppote"
	"uhppoted-mqtt/auth"
	"uhppoted-mqtt/mqtt"
	"uhppoted/config"
	"uhppoted/monitoring"
)

type Run struct {
	configuration       string
	dir                 string
	pidFile             string
	logFile             string
	logFileSize         int
	console             bool
	debug               bool
	healthCheckInterval time.Duration
	watchdogInterval    time.Duration
}

const (
	SERVICE = `uhppoted-mqtt`
)

func (c *Run) Name() string {
	return "run"
}

func (c *Run) Description() string {
	return fmt.Sprintf("Runs the %s daemon/service until terminated by the system service manager", SERVICE)
}

func (c *Run) Usage() string {
	return fmt.Sprintf("%s [--debug] [--config <file>] [--logfile <file>] [--logfilesize <bytes>] [--pid <file>]", SERVICE)
}

func (c *Run) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s <options>", SERVICE)
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
		log.Printf("WARN  Could not load configuration (%v)", err)
	}

	if err := os.MkdirAll(r.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("Unable to create working directory '%v': %v", r.dir, err)
	}

	pid := fmt.Sprintf("%d\n", os.Getpid())

	_, err := os.Stat(r.pidFile)
	if err == nil {
		return fmt.Errorf("PID lockfile '%v' already in use", r.pidFile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking PID lockfile '%v' (%v_)", r.pidFile, err)
	}

	if err := ioutil.WriteFile(r.pidFile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create PID lockfile: %v", err)
	}

	defer func() {
		os.Remove(r.pidFile)
	}()

	return f(conf)
}

func (r *Run) run(c *config.Config, logger *log.Logger, interrupt chan os.Signal) {
	logger.Printf("START")

	r.healthCheckInterval = c.HealthCheckInterval
	r.watchdogInterval = c.WatchdogInterval

	// ... initialise MQTT

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		ListenAddress:    c.ListenAddress,
		Devices:          make(map[uint32]*uhppote.Device),
		Debug:            r.debug,
	}

	for id, d := range c.Devices {
		u.Devices[id] = uhppote.NewDevice(id, d.Address, d.Rollover)
	}

	permissions, err := auth.NewPermissions(
		c.MQTT.Permissions.Enabled,
		c.MQTT.Permissions.Users,
		c.MQTT.Permissions.Groups,
		logger)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}

	mqttd := mqtt.MQTTD{
		ServerID: c.ServerID,
		TLS:      &tls.Config{},
		Connection: mqtt.Connection{
			Broker:   fmt.Sprintf(c.Connection.Broker),
			ClientID: c.Connection.ClientID,
			UserName: c.Connection.Username,
			Password: c.Connection.Password,
		},
		Topics: mqtt.Topics{
			Requests: c.Topics.Resolve(c.Topics.Requests),
			Replies:  c.Topics.Resolve(c.Topics.Replies),
			Events:   c.Topics.Resolve(c.Topics.Events),
			System:   c.Topics.Resolve(c.Topics.System),
		},
		Alerts: mqtt.Alerts{
			QOS:      c.Alerts.QOS,
			Retained: c.Alerts.Retained,
		},
		Encryption: mqtt.Encryption{
			SignOutgoing:    c.SignOutgoing,
			EncryptOutgoing: c.EncryptOutgoing,
			EventsKeyID:     c.EventsKeyID,
			SystemKeyID:     c.SystemKeyID,
			HOTP:            nil,
		},
		Authentication: c.Authentication,
		Permissions:    *permissions,
		EventMap:       c.EventIDs,
		Debug:          r.debug,
	}

	// ... TLS

	if strings.HasPrefix(mqttd.Connection.Broker, "tls:") {
		pem, err := ioutil.ReadFile(c.Connection.BrokerCertificate)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		} else {
			mqttd.TLS.InsecureSkipVerify = false
			mqttd.TLS.RootCAs = x509.NewCertPool()

			if ok := mqttd.TLS.RootCAs.AppendCertsFromPEM(pem); !ok {
				logger.Printf("ERROR: Could not initialise MQTTD CA certificates")
			}
		}

		certificate, err := tls.LoadX509KeyPair(c.Connection.ClientCertificate, c.Connection.ClientKey)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		} else {
			mqttd.TLS.Certificates = []tls.Certificate{certificate}
		}
	}

	// ... authentication

	hmac, err := auth.NewHMAC(c.HMAC.Required, c.HMAC.Key)
	if c.HMAC.Required && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	hotp, err := auth.NewHOTP(c.MQTT.HOTP.Range, c.MQTT.HOTP.Secrets, c.MQTT.HOTP.Counters, logger)
	if mqttd.Authentication == "HOTP" && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	rsa, err := auth.NewRSA(c.RSA.KeyDir, logger)
	if mqttd.Authentication == "RSA" && err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	nonce, err := auth.NewNonce(c.Nonce.Required, c.Nonce.Server, c.Nonce.Clients, logger)
	if err != nil {
		logger.Printf("ERROR: %v", err)
		return
	}

	mqttd.HMAC = *hmac
	mqttd.Encryption.HOTP = hotp
	mqttd.Encryption.RSA = rsa
	mqttd.Encryption.Nonce = *nonce

	// ... monitoring

	healthcheck := monitoring.NewHealthCheck(&u, c.HealthCheckIdle, c.HealthCheckIgnore, logger)

	// ... listen

	err = r.listen(&u, &mqttd, &healthcheck, logger, interrupt)
	if err != nil {
		logger.Printf("ERROR %v", err)
	}

	logger.Printf("INFO  exit")
}

func (r *Run) listen(
	u *uhppote.UHPPOTE,
	mqttd *mqtt.MQTTD,
	healthcheck *monitoring.HealthCheck,
	logger *log.Logger,
	interrupt chan os.Signal) error {

	// ... MQTT

	pid := fmt.Sprintf("%d\n", os.Getpid())
	workdir := filepath.Dir(r.pidFile)
	lockfile := filepath.Join(workdir, fmt.Sprintf("%s.lock", mqttd.Connection.ClientID))

	_, err := os.Stat(lockfile)
	if err == nil {
		return fmt.Errorf("MQTT client lockfile '%v' already in use", lockfile)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking MQTT client lockfile '%v' (%v_)", lockfile, err)
	}

	if err := ioutil.WriteFile(lockfile, []byte(pid), 0644); err != nil {
		return fmt.Errorf("Unable to create MQTT client lockfile: %v", err)
	}

	defer func() {
		os.Remove(lockfile)
	}()

	if err := mqttd.Run(u, logger); err != nil {
		return err
	}

	defer mqttd.Close(logger)

	// ... monitoring

	monitor := mqtt.NewSystemMonitor(mqttd, logger)
	watchdog := monitoring.NewWatchdog(healthcheck, logger)
	k := time.NewTicker(r.healthCheckInterval)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck.Exec(monitor)
		}
	}()

	// ... wait until interrupted

	w := time.NewTicker(r.watchdogInterval)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog.Exec(monitor); err != nil {
				return err
			}

		case <-interrupt:
			logger.Printf("... interrupt")
			return nil
		}
	}
}
