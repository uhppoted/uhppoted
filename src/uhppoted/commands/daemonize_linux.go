package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Daemonize struct {
}

const service = `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}
User=uhppoted
Group=uhppoted

[Install]
WantedBy=multi-user.target
`

const rotate = `{{range .}}{{.LogFile}} {
    daily
    rotate 30
    compress
        compresscmd /bin/bzip2
        compressext .bz2
        dateext
    missingok
    notifempty
    su uhppoted uhppoted
    postrotate
       /usr/bin/killall -HUP uhppoted
    endscript
}{{end}}
`

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	if err := c.systemd(); err != nil {
		return err
	}

	if err := c.logrotate(); err != nil {
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

func (c *Daemonize) logrotate() error {
	data := []struct {
		LogFile string
	}{
		{LogFile: "/var/log/uhppoted/uhppoted.log"},
	}

	path := filepath.Join("/etc/logrotate.d", "uhppoted")
	t := template.Must(template.New("uhppoted.logrotate").Parse(rotate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Daemonize) mkdirs() error {
	directories := []string{
		"/var/uhppoted",
		"/var/log/uhppoted",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		err := os.MkdirAll(dir, 0774)
		if err != nil {
			return err
		}

		cmd := exec.Command("chown", "-R", "uhppoted:uhppoted", dir)
		out, err := cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to set ownership of '%s' to uhppoted:uhppoted (%v)\n", dir, err)
			return err
		}
	}

	return nil
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
