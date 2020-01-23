package commands

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"uhppote"
	"uhppoted-mqtt/auth"
	"uhppoted-mqtt/config"
	"uhppoted-mqtt/mqtt"
	"uhppoted/monitoring"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	console       bool
	debug         bool
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
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
		log.Printf("\n   WARN  Could not load configuration (%v)\n\n", err)
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

	// ... initialise MQTT

	u := uhppote.UHPPOTE{
		BindAddress:      c.BindAddress,
		BroadcastAddress: c.BroadcastAddress,
		ListenAddress:    c.ListenAddress,
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            r.debug,
	}

	for id, d := range c.Devices {
		if d.Address != nil {
			u.Devices[id] = d.Address
		}
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
		Broker:   fmt.Sprintf(c.Broker),
		TLS:      &tls.Config{},
		Topics: mqtt.Topics{
			Requests: c.Topics.Resolve(c.Topics.Requests),
			Replies:  c.Topics.Resolve(c.Topics.Replies),
			Events:   c.Topics.Resolve(c.Topics.Events),
			System:   c.Topics.Resolve(c.Topics.System),
		},
		Encryption: mqtt.Encryption{
			SignOutgoing:    c.SignOutgoing,
			EncryptOutgoing: c.EncryptOutgoing,
			EventsKeyID:     c.EventsKeyID,
			SystemKeyID:     c.SystemKeyID,
			HOTP:            nil,
		},
		Authentication:      c.Authentication,
		Permissions:         *permissions,
		EventMap:            c.EventIDs,
		HealthCheckInterval: c.MQTT.HealthCheckInterval,
		Debug:               r.debug,
	}

	// ... TLS

	if strings.HasPrefix(mqttd.Broker, "tls:") {
		pem, err := ioutil.ReadFile(c.BrokerCertificate)
		if err != nil {
			logger.Printf("ERROR: %v", err)
		} else {
			mqttd.TLS.InsecureSkipVerify = false
			mqttd.TLS.RootCAs = x509.NewCertPool()

			if ok := mqttd.TLS.RootCAs.AppendCertsFromPEM(pem); !ok {
				logger.Printf("ERROR: Could not initialise MQTTD CA certificates")
			}
		}

		certificate, err := tls.LoadX509KeyPair(c.ClientCertificate, c.ClientKey)
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

	// ... listen forever

	for {
		err := r.listen(&u, &mqttd, logger, interrupt)
		if err != nil {
			logger.Printf("ERROR %v", err)
			continue
		}

		logger.Printf("INFO  exit\n")
		break
	}

	logger.Printf("STOP")
}

func (r *Run) listen(u *uhppote.UHPPOTE, mqttd *mqtt.MQTTD, logger *log.Logger, interrupt chan os.Signal) error {
	// ... MQTT task

	go func() {
		mqttd.Run(u, logger)
	}()

	defer mqttd.Close(logger)

	// ... health-check task

	healthcheck := monitoring.NewHealthCheck(u, logger)
	watchdog := monitoring.NewWatchdog(&healthcheck, logger)
	k := time.NewTicker(mqttd.HealthCheckInterval)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck.Exec(mqttd)
		}
	}()

	// ... wait until interrupted/closed

	closed := make(chan struct{})
	w := time.NewTicker(5 * time.Second)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog.Exec(mqttd); err != nil {
				return err
			}

		case <-interrupt:
			logger.Printf("... interrupt")
			return nil

		case <-closed:
			logger.Printf("... closed")
			return errors.New("MQTT client error")
		}
	}
}
