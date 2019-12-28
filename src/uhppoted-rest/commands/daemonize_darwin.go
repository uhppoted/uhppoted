package commands

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"uhppoted-rest/config"
	xpath "uhppoted-rest/encoding/plist"
)

type info struct {
	Label            string
	Executable       string
	ConfigDirectory  string
	WorkingDirectory string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

type plist struct {
	Label             string
	Program           string
	WorkingDirectory  string
	ProgramArguments  []string
	KeepAlive         bool
	RunAtLoad         bool
	StandardOutPath   string
	StandardErrorPath string
}

const newsyslog = `#logfilename                                       [owner:group]  mode  count  size   when  flags [/pid_file]  [sig_num]
{{range .}}{{.LogFile}}  :              644   30     10000  @T00  J     {{.PID}}
{{end}}`

const confTemplate = `# UDP
bind.address = {{.BindAddress}}
broadcast.address = {{.BroadcastAddress}}

# REST API
rest.http.enabled = false
rest.http.port = 8080
rest.https.enabled = true
rest.https.port = 8443
rest.tls.key = {{.ConfigDirectory}}/rest/uhppoted.key
rest.tls.certificate = {{.ConfigDirectory}}/rest/uhppoted.cert
rest.tls.ca = {{.ConfigDirectory}}/rest/ca.cert

# OPEN API
# openapi.enabled = false
# openapi.directory = {{.WorkDir}}\rest\openapi

# DEVICES
# Example configuration for UTO311-L04 with serial number 405419896
# UT0311-L0x.405419896.address = 192.168.1.100:60000
# UT0311-L0x.405419896.door.1 = Front Door
# UT0311-L0x.405419896.door.2 = Side Door
# UT0311-L0x.405419896.door.3 = Garage
# UT0311-L0x.405419896.door.4 = Workshop
`

type Daemonize struct {
}

func NewDaemonize() *Daemonize {
	return &Daemonize{}
}

func (c *Daemonize) Name() string {
	return "daemonize"
}

func (c *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (c *Daemonize) Description() string {
	return "Daemonizes uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println()
	fmt.Println("  Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println("    Daemonizes uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}

func (c *Daemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... daemonizing")
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast := config.DefaultIpAddresses()

	d := info{
		Label:            "com.github.twystd.uhppoted-rest",
		Executable:       executable,
		ConfigDirectory:  "/usr/local/etc/com.github.twystd.uhppoted",
		WorkingDirectory: "/usr/local/var/com.github.twystd.uhppoted",
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := c.launchd(&d); err != nil {
		return err
	}

	if err := c.mkdirs(); err != nil {
		return err
	}

	if err := c.logrotate(); err != nil {
		return err
	}

	if err := c.firewall(); err != nil {
		return err
	}

	if err := c.conf(&d); err != nil {
		return err
	}

	fmt.Println("   ... com.github.twystd.uhppoted-rest registered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("   sudo launchctl load /Library/LaunchDaemons/com.github.twystd.uhppoted-rest.plist")
	fmt.Println()

	return nil
}

func (c *Daemonize) launchd(d *info) error {
	path := filepath.Join("/Library/LaunchDaemons", "com.github.twystd.uhppoted-rest.plist")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	pl := plist{
		Label:             d.Label,
		Program:           d.Executable,
		WorkingDirectory:  "/usr/local/var/com.github.twystd.uhppoted",
		ProgramArguments:  []string{},
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   "/usr/local/var/log/com.github.twystd.uhppoted-rest.log",
		StandardErrorPath: "/usr/local/var/log/com.github.twystd.uhppoted-rest.err",
	}

	if !os.IsNotExist(err) {
		current, err := c.parse(path)
		if err != nil {
			return err
		}

		pl.WorkingDirectory = current.WorkingDirectory
		pl.ProgramArguments = current.ProgramArguments
		pl.KeepAlive = current.KeepAlive
		pl.RunAtLoad = current.RunAtLoad
		pl.StandardOutPath = current.StandardOutPath
		pl.StandardErrorPath = current.StandardErrorPath
	}

	return c.daemonize(path, pl)
}

func (c *Daemonize) parse(path string) (*plist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := plist{}
	decoder := xpath.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (c *Daemonize) daemonize(path string, p interface{}) error {
	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := xpath.NewEncoder(f)
	if err = encoder.Encode(p); err != nil {
		return err
	}

	return nil
}

func (c *Daemonize) mkdirs() error {
	dir := "/usr/local/var/com.github.twystd.uhppoted"

	fmt.Printf("   ... creating '%s'\n", dir)

	return os.MkdirAll(dir, 0644)
}

func (c *Daemonize) conf(d *info) error {
	path := filepath.Join(d.ConfigDirectory, "uhppoted.conf")
	t := template.Must(template.New("uhppoted.conf").Parse(confTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (c *Daemonize) logrotate() error {
	dir := "/usr/local/var/log"
	pid := "/usr/local/var/com.github.twystd.uhppoted/uhppoted-rest.pid"
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: filepath.Join(dir, "com.github.twystd.uhppoted-rest.log"),
			PID:     pid,
		},
		{
			LogFile: filepath.Join(dir, "com.github.twystd.uhppoted-rest.err"),
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", "uhppoted-rest.conf")

	fmt.Printf("   ... creating '%s'\n", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, logfiles)
}

func (c *Daemonize) firewall() error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Println("   *** WARNING: adding 'uhppoted-rest' to the application firewall and unblocking incoming connections")
	fmt.Println("   ***")
	fmt.Println()

	path, err := os.Executable()
	if err != nil {
		return fmt.Errorf("Failed to get path to executable: %v\n", err)
	}

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to retrieve application firewall global state (%v)\n", err)
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to disable the application firewall (%v)\n", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to add 'uhppoted-rest' to the application firewall (%v)\n", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to unblock 'uhppoted-rest' on the application firewall (%v)\n", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to re-enable the application firewall (%v)\n", err)
		}

		fmt.Println()
	}

	return nil
}
