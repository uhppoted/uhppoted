package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"uhppoted/encoding/plist"
)

type info struct {
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

type Daemonize struct {
}

func (c *Daemonize) Execute(ctx Context) error {
	fmt.Println("   ... daemonizing")

	if err := c.launchd(); err != nil {
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

	fmt.Println("   ... com.github.twystd.uhppoted registered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("   sudo launchctl load /Library/LaunchDaemons/com.github.twystd.uhppoted")
	fmt.Println()

	return nil
}

func (c *Daemonize) launchd() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	plist := struct {
		Label             string
		Program           string
		WorkingDirectory  string
		ProgramArguments  []string
		KeepAlive         bool
		RunAtLoad         bool
		StandardOutPath   string
		StandardErrorPath string
	}{
		Label:             "com.github.twystd.uhppoted",
		Program:           executable,
		WorkingDirectory:  "/usr/local/var/com.github.twystd.uhppoted",
		ProgramArguments:  []string{},
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   "/usr/local/var/log/com.github.twystd.uhppoted.log",
		StandardErrorPath: "/usr/local/var/log/com.github.twystd.uhppoted.err",
	}

	path := filepath.Join("/Library/LaunchDaemons", "com.github.twystd.uhppoted.plist")
	_, err = os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !os.IsNotExist(err) {
		current, err := c.parse(path)
		if err != nil {
			return err
		}

		plist.WorkingDirectory = current.WorkingDirectory
		plist.ProgramArguments = current.ProgramArguments
		plist.KeepAlive = current.KeepAlive
		plist.RunAtLoad = current.RunAtLoad
		plist.StandardOutPath = current.StandardOutPath
		plist.StandardErrorPath = current.StandardErrorPath
	}

	return c.daemonize(path, plist)
}

func (c *Daemonize) parse(path string) (*info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := info{}
	decoder := plist.NewDecoder(f)
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

	encoder := plist.NewEncoder(f)
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

func (c *Daemonize) logrotate() error {
	dir := "/usr/local/var/log"
	pid := "/usr/local/var/com.github.twystd.uhppoted/uhppoted.pid"
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: filepath.Join(dir, "com.github.twystd.uhppoted.log"),
			PID:     pid,
		},
		{
			LogFile: filepath.Join(dir, "com.github.twystd.uhppoted.err"),
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", "uhppoted.conf")

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
	fmt.Println("   *** WARNING: adding 'uhppoted' to the application firewall and unblocking incoming connections")
	fmt.Println("   ***")
	fmt.Println()

	path, err := os.Executable()
	if err != nil {
		fmt.Errorf("Failed to get path to executable: %v\n", err)
		return err
	}

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		fmt.Errorf("ERROR: Failed to retrieve application firewall global state (%v)\n", err)
		return err
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to disable the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to add 'uhppoted' to the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to unblock 'uhppoted' on the application firewall (%v)\n", err)
			return err
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			fmt.Errorf("ERROR: Failed to re-enable the application firewall (%v)\n", err)
			return err
		}

		fmt.Println()
	}

	return nil
}

func (c *Daemonize) Cmd() string {
	return "daemonize"
}

func (c *Daemonize) Description() string {
	return "Daemonizes uhppoted as a service/daemon"
}

func (c *Daemonize) Usage() string {
	return ""
}

func (c *Daemonize) Help() {
	fmt.Println("Usage: uhppoted daemonize")
	fmt.Println()
	fmt.Println(" Daemonizes uhppoted as a service/daemon that runs on startup")
	fmt.Println()
}
