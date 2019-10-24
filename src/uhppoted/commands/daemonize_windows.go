package commands

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"net"
	"os"
	"path/filepath"
	"text/template"
	"uhppoted/config"
)

type Daemonize struct {
	name        string
	description string
}

type data struct {
	Name             string
	Description      string
	Executable       string
	WorkDir          string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

const confTemplate = `bind.address = {{.BindAddress}}
broadcast.address = {{.BroadcastAddress}}

# Example configuration for UTO311-L04 with serial number 305419896
# UT0311-L0x.305419896.address = 192.168.1.100:60000
# UT0311-L0x.305419896.door.1 = Front Door
# UT0311-L0x.305419896.door.2 = Side Door
# UT0311-L0x.305419896.door.3 = Garage
# UT0311-L0x.305419896.door.4 = Workshop
`

func NewDaemonize() *Daemonize {
	return &Daemonize{
		name: "uhppoted",
	}
}

func (c *Daemonize) Parse(args []string) error {
	return nil
}

func (c *Daemonize) Execute(ctx Context) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast, err := config.DefaultIpAddresses()
	if err != nil {
		return err
	}

	if bind == nil || broadcast == nil {
		return errors.New("Unable to determine default bind and broadcast IP addresses")
	}

	d := data{
		Name:             c.name,
		Description:      "UHPPOTE UTO311-L0x access card controllers service/daemon ",
		Executable:       executable,
		WorkDir:          workdir(),
		BindAddress:      bind,
		BroadcastAddress: broadcast,
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

	fmt.Println("   ... uhppoted registered as a Windows system service")
	fmt.Println()
	fmt.Println("   The service will start automatically on the next system restart. Start it manually from the")
	fmt.Println("   'Services' application or from the command line by executing the following command:")
	fmt.Println()
	fmt.Println("     > net start uhppoted")
	fmt.Println()

	return nil
}

func (c *Daemonize) register(d *data) error {
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

func (c *Daemonize) mkdirs(d *data) error {
	directories := []string{
		d.WorkDir,
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}
	}

	return nil
}

func (c *Daemonize) conf(d *data) error {
	path := filepath.Join(d.WorkDir, "uhppoted.conf")
	t := template.Must(template.New("uhppoted.conf").Parse(confTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func workdir() string {
	programData, err := windows.KnownFolderPath(windows.FOLDERID_ProgramData, windows.KF_FLAG_DEFAULT)
	if err != nil {
		return `C:\uhppoted`
	}

	return filepath.Join(programData, "uhppoted")
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted as a Windows service"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Registers uhppoted as a windows Service that runs on startup")
	fmt.Println()
}
