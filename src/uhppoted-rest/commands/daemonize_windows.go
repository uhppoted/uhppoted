package commands

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"uhppoted-rest/config"
)

type Daemonize struct {
	name        string
	description string
}

type info struct {
	Name             string
	Description      string
	Executable       string
	WorkDir          string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

const confTemplate = `# UDP
bind.address = {{.BindAddress}}
broadcast.address = {{.BroadcastAddress}}

# REST API
rest.http.enabled = false
rest.http.port = 8080
rest.https.enabled = true
rest.https.port = 8443
rest.tls.key = {{.WorkDir}}\rest\uhppoted.key
rest.tls.certificate = {{.WorkDir}}\rest\uhppoted.cert
rest.tls.ca = {{.WorkDir}}\rest\ca.cert
# rest.openapi.enabled = false
# rest.openapi.directory = {{.WorkDir}}\rest\openapi

# DEVICES
# Example configuration for UTO311-L04 with serial number 305419896
# UT0311-L0x.305419896.address = 192.168.1.100:60000
# UT0311-L0x.305419896.door.1 = Front Door
# UT0311-L0x.305419896.door.2 = Side Door
# UT0311-L0x.305419896.door.3 = Garage
# UT0311-L0x.305419896.door.4 = Workshop
`

func NewDaemonize() *Daemonize {
	return &Daemonize{
		name:        "uhppoted-rest",
		description: "UHPPOTE UTO311-L0x access card controllers service/daemon",
	}
}

func (c *Daemonize) Name() string {
	return "daemonize"
}

func (c *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted-rest as a Windows service"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted-rest daemonize")
	fmt.Println()
	fmt.Println("    Registers uhppoted-rest as a windows Service that runs on startup")
	fmt.Println()
}
func (c *Daemonize) Execute(ctx context.Context) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast := config.DefaultIpAddresses()

	d := info{
		Name:             c.name,
		Description:      c.description,
		Executable:       executable,
		WorkDir:          workdir(),
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := c.register(&d); err != nil {
		return err
	}

	if err := c.mkdirs(&d); err != nil {
		return err
	}

	if err := c.conf(&d); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted-rest registered as a Windows system service")
	fmt.Println()
	fmt.Println("   The service will start automatically on the next system restart. Start it manually from the")
	fmt.Println("   'Services' application or from the command line by executing the following command:")
	fmt.Println()
	fmt.Println("     > net start uhppoted-rest")
	fmt.Println()

	return nil
}

func (c *Daemonize) register(d *info) error {
	config := mgr.Config{
		DisplayName:      d.Name,
		Description:      d.Description,
		StartType:        mgr.StartAutomatic,
		DelayedAutoStart: true,
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(d.Name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", d.Name)
	}

	s, err = m.CreateService(d.Name, d.Executable, config, "is", "auto-started")
	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(d.Name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("InstallAsEventCreate() failed: %v", err)
	}

	return nil
}

func (c *Daemonize) mkdirs(d *info) error {
	directories := []string{
		d.WorkDir,
		filepath.Join(d.WorkDir, "rest"),
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}
	}

	return nil
}

func (c *Daemonize) conf(d *info) error {
	path := filepath.Join(d.WorkDir, "uhppoted.conf")
	t := template.Must(template.New("uhppoted.conf").Parse(confTemplate))
	var b strings.Builder

	fmt.Printf("   ... creating '%s'\n", path)

	if err := t.Execute(&b, d); err != nil {
		return err
	}

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
