package commands

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"net"
	"os"
	"path/filepath"
	"strings"
	"uhppoted/config"
)

var DAEMONIZE = Daemonize{
	name:        SERVICE,
	description: "UHPPOTE UTO311-L0x access card controllers service",
	workdir:     workdir(),
	logdir:      filepath.Join(workdir(), "logs"),
	config:      filepath.Join(workdir(), "uhppoted.conf"),
	hotp:        filepath.Join(workdir(), "mqtt.hotp.secrets"),
}

type info struct {
	Executable       string
	WorkDir          string
	LogDir           string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

type Daemonize struct {
	name        string
	description string
	workdir     string
	logdir      string
	config      string
	hotp        string
}

func (d *Daemonize) Name() string {
	return "daemonize"
}

func (d *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (d *Daemonize) Description() string {
	return fmt.Sprintf("Registers %s as a Windows service", SERVICE)
}

func (d *Daemonize) Usage() string {
	return ""
}

func (d *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a Windows service\n", SERVICE)
	fmt.Println()
}

func (d *Daemonize) Execute(ctx context.Context) error {
	dir := filepath.Dir(d.config)
	r := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Printf("     **** PLEASE MAKE SURE YOU HAVE A BACKUP COPY OF THE CONFIGURATION INFORMATION AND KEYS IN %s ***\n", dir)
	fmt.Println()
	fmt.Printf("     Enter 'yes' to continue with the installation: ")

	text, err := r.ReadString('\n')
	if err != nil || strings.TrimSpace(text) != "yes" {
		fmt.Println()
		fmt.Printf("     -- installation cancelled --")
		fmt.Println()
		return nil
	}

	return d.execute(ctx)
}

func (d *Daemonize) execute(ctx context.Context) error {
	fmt.Println()
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast, _ := config.DefaultIpAddresses()

	i := info{
		Executable:       executable,
		WorkDir:          d.workdir,
		LogDir:           d.logdir,
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := d.register(&i); err != nil {
		return err
	}

	if err := d.mkdirs(&i); err != nil {
		return err
	}

	if err := d.conf(&i); err != nil {
		return err
	}

	if err := d.genkeys(&i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a Windows system service\n", SERVICE)
	fmt.Println()
	fmt.Println("   The service will start automatically on the next system restart. Start it manually from the")
	fmt.Println("   'Services' application or from the command line by executing the following command:")
	fmt.Println()
	fmt.Printf("     > net start %s\n", SERVICE)
	fmt.Printf("     > sc query %s\n", SERVICE)
	fmt.Println()
	fmt.Println("   Please replace the default RSA keys for event and system messages:")
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(d.config), "mqtt", "rsa", "encryption", "event.pub"))
	fmt.Printf("     - %s\n", filepath.Join(filepath.Dir(d.config), "mqtt", "rsa", "encryption", "system.pub"))
	fmt.Println()

	return nil
}

func (d *Daemonize) register(i *info) error {
	config := mgr.Config{
		DisplayName:      d.name,
		Description:      d.description,
		StartType:        mgr.StartAutomatic,
		DelayedAutoStart: true,
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(d.name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", d.Name)
	}

	s, err = m.CreateService(d.name, i.Executable, config, "is", "auto-started")
	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(d.name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("InstallAsEventCreate() failed: %v", err)
	}

	return nil
}

func (d *Daemonize) mkdirs(i *info) error {
	directories := []string{
		i.WorkDir,
		i.LogDir,
		filepath.Join(i.WorkDir, "mqtt"),
		filepath.Join(i.WorkDir, "mqtt", "rsa"),
		filepath.Join(i.WorkDir, "mqtt", "rsa", "encryption"),
		filepath.Join(i.WorkDir, "mqtt", "rsa", "signing"),
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemonize) conf(i *info) error {
	path := d.config

	fmt.Printf("   ... creating '%s'\n", path)

	// initialise config from existing uhppoted.conf
	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	// generate HMAC and RSA keys
	if cfg.MQTT.HMAC.Key == "" {
		hmac, err := hmac()
		if err != nil {
			return err
		}

		cfg.MQTT.HMAC.Key = hmac
	}

	// replace line endings
	var b strings.Builder

	err := cfg.Write(&b)

	// write back config with any updated information
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	replacer := strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
	)

	if _, err = f.Write([]byte(replacer.Replace(b.String()))); err != nil {
		return err
	}

	return nil
}

func (d *Daemonize) genkeys(i *info) error {
	return genkeys(filepath.Dir(d.config), d.hotp)
}
