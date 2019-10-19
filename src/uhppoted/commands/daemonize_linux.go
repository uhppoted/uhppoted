package commands

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"text/template"
)

type Daemonize struct {
}

type data struct {
	Description   string
	Documentation string
	Executable    string
	PID           string
	User          string
	Group         string
	Uid           int
	Gid           int
	LogFiles      []string
}

var templates = map[string]string{
	"service": `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`,

	"rotate": `{{range .LogFiles}}{{.}} {
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
`,
}

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	uid, gid, err := c.getUser()
	if err != nil {
		return err
	}

	d := data{
		Description:   "UHPPOTE UTO311-L0x access card controllers service/daemon ",
		Documentation: "https://github.com/twystd/uhppote-go",
		Executable:    executable,
		PID:           "/var/uhppoted/uhppoted.pid",
		User:          "uhppoted",
		Group:         "uhppoted",
		Uid:           uid,
		Gid:           gid,
		LogFiles:      []string{"/var/log/uhppoted/uhppoted.log"},
	}

	if err := c.systemd(&d); err != nil {
		return err
	}

	if err := c.logrotate(&d); err != nil {
		return err
	}

	if err := c.mkdirs(&d); err != nil {
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

func (c *Daemonize) getUser() (int, int, error) {
	u, err := user.Lookup("uhppoted")
	if err != nil {
		return 0, 0, err
	}

	g, err := user.LookupGroup("uhppoted")
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

func (c *Daemonize) systemd(d *data) error {
	path := filepath.Join("/etc/systemd/system", "uhppoted.service")
	t := template.Must(template.New("uhppoted.service").Parse(templates["service"]))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (c *Daemonize) logrotate(d *data) error {
	path := filepath.Join("/etc/logrotate.d", "uhppoted")
	t := template.Must(template.New("uhppoted.logrotate").Parse(templates["rotate"]))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (c *Daemonize) mkdirs(d *data) error {
	directories := []string{
		"/var/uhppoted",
		"/var/log/uhppoted",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}

		if err := os.Chown(dir, d.Uid, d.Gid); err != nil {
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
