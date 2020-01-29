package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	xpath "uhppoted-rest/encoding/plist"
	"uhppoted/config"
)

type info struct {
	Label      string
	Executable string
	WorkDir    string
	StdLogFile string
	ErrLogFile string
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

// const confTemplate = `# UDP
// bind.address = {{.BindAddress}}
// broadcast.address = {{.BroadcastAddress}}
//
// # REST API
// rest.http.enabled = false
// rest.http.port = 8080
// rest.https.enabled = true
// rest.https.port = 8443
// rest.tls.key = {{.ConfigDirectory}}/rest/uhppoted.key
// rest.tls.certificate = {{.ConfigDirectory}}/rest/uhppoted.cert
// rest.tls.ca = {{.ConfigDirectory}}/rest/ca.cert
//
// # OPEN API
// # openapi.enabled = false
// # openapi.directory = {{.WorkDir}}\rest\openapi
//
// # DEVICES
// # Example configuration for UTO311-L04 with serial number 405419896
// # UT0311-L0x.405419896.address = 192.168.1.100:60000
// # UT0311-L0x.405419896.door.1 = Front Door
// # UT0311-L0x.405419896.door.2 = Side Door
// # UT0311-L0x.405419896.door.3 = Garage
// # UT0311-L0x.405419896.door.4 = Workshop
// `

var DAEMONIZE = Daemonize{
	plist:   fmt.Sprintf("com.github.twystd.%s.plist", SERVICE),
	workdir: "/usr/local/var/com.github.twystd.uhppoted",
	logdir:  "/usr/local/var/com.github.twystd.uhppoted/log",
	config:  fmt.Sprintf("/usr/local/etc/com.github.twystd.uhppoted/%s.stuff", SERVICE),
}

type Daemonize struct {
	plist   string
	workdir string
	logdir  string
	config  string
}

func NewDaemonize() *Daemonize {
	return &Daemonize{}
}

func (d *Daemonize) Name() string {
	return "daemonize"
}

func (d *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (d *Daemonize) Description() string {
	return fmt.Sprintf("Daemonizes %s as a service/daemon", SERVICE)
}

func (d *Daemonize) Usage() string {
	return ""
}

func (d *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Daemonizes %s as a service/daemon that runs on startup\n", SERVICE)
	fmt.Println()
}

func (d *Daemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	i := info{
		Label:      fmt.Sprintf("com.github.twystd.%s", SERVICE),
		Executable: executable,
		WorkDir:    d.workdir,
		StdLogFile: filepath.Join(d.logdir, fmt.Sprintf("%s.log", SERVICE)),
		ErrLogFile: filepath.Join(d.logdir, fmt.Sprintf("%s.err", SERVICE)),
	}

	if err := d.launchd(&i); err != nil {
		return err
	}

	if err := d.mkdirs(); err != nil {
		return err
	}

	if err := d.logrotate(&i); err != nil {
		return err
	}

	if err := d.firewall(&i); err != nil {
		return err
	}

	if err := d.conf(&i); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a LaunchDaemon\n", i.Label)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf("   sudo launchctl load /Library/LaunchDaemons/com.github.twystd.%s.plist\n", SERVICE)
	fmt.Println()

	return nil
}

func (d *Daemonize) launchd(i *info) error {
	path := filepath.Join("/Library/LaunchDaemons", d.plist)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	pl := plist{
		Label:             i.Label,
		Program:           i.Executable,
		WorkingDirectory:  i.WorkDir,
		ProgramArguments:  []string{},
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   i.StdLogFile,
		StandardErrorPath: i.ErrLogFile,
	}

	if !os.IsNotExist(err) {
		current, err := d.parse(path)
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

	return d.daemonize(path, pl)
}

func (d *Daemonize) parse(path string) (*plist, error) {
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

func (d *Daemonize) daemonize(path string, p interface{}) error {
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

func (d *Daemonize) mkdirs() error {
	directories := []string{
		d.workdir,
		d.logdir,
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemonize) conf(i *info) error {
	path := d.config

	fmt.Printf("   ... creating '%s'\n", path)

	cfg := config.NewConfig()

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	cfg.Write(f)

	return nil
}

func (c *Daemonize) logrotate(i *info) error {
	pid := filepath.Join(c.workdir, fmt.Sprintf("%s.pid", SERVICE))
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: i.StdLogFile,
			PID:     pid,
		},
		{
			LogFile: i.ErrLogFile,
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", fmt.Sprintf("%s.conf", SERVICE))

	fmt.Printf("   ... creating '%s'\n", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, logfiles)
}

func (c *Daemonize) firewall(i *info) error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Printf("   *** WARNING: adding '%s' to the application firewall and unblocking incoming connections\n", SERVICE)
	fmt.Println("   ***")
	fmt.Println()

	path := i.Executable

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("Failed to retrieve application firewall global state (%v)", err)
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to disable the application firewall (%v)", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to add 'uhppoted-rest' to the application firewall (%v)", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to unblock 'uhppoted-rest' on the application firewall (%v)", err)
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to re-enable the application firewall (%v)", err)
		}

		fmt.Println()
	}

	return nil
}
