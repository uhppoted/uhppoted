package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type Daemonize struct {
}

const service = `
[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}

[Install]
WantedBy=multi-user.target
`

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	if err := c.systemd(); err != nil {
		return err
	}

	if err := c.mkdirs(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted registered as a systemd service")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("   sudo systemctl start uhppoted")
	fmt.Println()

	return nil
}

func (c *Daemonize) systemd() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	data := struct {
		Description   string
		Documentation string
		Executable    string
		PID           string
	}{
		Description:   "UHPPOTE UTO311-L0x access card controllers service/daemon ",
		Documentation: "https://github.com/twystd/uhppote-go",
		Executable:    executable,
		PID:           "/var/uhppoted/uhppoted.pid",
	}

	path := filepath.Join("/etc/systemd/system", "uhppoted.service")
	t := template.Must(template.New("uhppoted.service").Parse(service))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, data)
}

func (c *Daemonize) mkdirs() error {
	dir := "/var/uhppoted"

	fmt.Printf("   ... creating '%s'\n", dir)

	return os.MkdirAll(dir, 0644)
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Registers uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Registers uhppoted as a systemd service/daemon that runs on startup")
	fmt.Println()
}
