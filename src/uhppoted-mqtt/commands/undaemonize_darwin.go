package commands

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	xpath "uhppoted-rest/encoding/plist"
)

var UNDAEMONIZE = Undaemonize{
	plist:   fmt.Sprintf("com.github.twystd.%s.plist", SERVICE),
	workdir: "/usr/local/var/com.github.twystd.uhppoted",
	logdir:  "/usr/local/var/com.github.twystd.uhppoted/log",
	config:  fmt.Sprintf("/usr/local/etc/com.github.twystd.uhppoted/%s.stuff", SERVICE),
}

type Undaemonize struct {
	plist   string
	workdir string
	logdir  string
	config  string
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{}
}

func (u *Undaemonize) Name() string {
	return "undaemonize"
}

func (u *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (u *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters %s as a service/daemon", SERVICE)
}

func (u *Undaemonize) Usage() string {
	return ""
}

func (u *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s from launchd as a service/daemon", SERVICE)
	fmt.Println()
}

func (u *Undaemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... undaemonizing")

	executable, err := u.launchd()
	if err != nil {
		return err
	}

	if err := u.logrotate(); err != nil {
		return err
	}

	if err := u.clean(); err != nil {
		return err
	}

	if err := u.firewall(executable); err != nil {
		return err
	}

	fmt.Printf("   ... com.github.twystd.%s unregistered as a LaunchDaemon\n", SERVICE)
	fmt.Println()

	return nil
}

func (u *Undaemonize) parse(path string) (*info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	i := info{}
	decoder := xpath.NewDecoder(f)
	err = decoder.Decode(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (u *Undaemonize) launchd() (string, error) {
	label := fmt.Sprintf("com.github.twystd.%s", SERVICE)

	path := filepath.Join("/Library/LaunchDaemons", u.plist)
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... nothing to do for 'launchd'   (%s does not exist)\n", path)
		return "", nil
	}

	// stop daemon
	fmt.Println("   ... unloading LaunchDaemon")
	cmd := exec.Command("launchctl", "unload", path)
	out, err := cmd.CombinedOutput()
	fmt.Println()
	fmt.Printf("   > %s", out)
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("Failed to unload '%s' (%v)\n", label, err)
	}

	// get launchd executable from plist

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	pl := plist{}
	decoder := xpath.NewDecoder(f)
	if err = decoder.Decode(&pl); err != nil {
		f.Close()
		return "", err
	}

	f.Close()

	// remove plist file
	fmt.Printf("   ... removing '%s'\n", path)
	err = os.Remove(path)
	if err != nil {
		return pl.Program, err
	}

	return pl.Program, nil
}

func (u *Undaemonize) logrotate() error {
	path := filepath.Join("/etc/newsyslog.d", fmt.Sprintf("%s.conf", SERVICE))

	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... nothing to do for 'newsyslog' (%s does not exist)\n", path)
		return nil
	}

	fmt.Printf("   ... removing '%s'\n", path)

	return os.Remove(path)
}

func (u *Undaemonize) clean() error {
	files := []string{
		filepath.Join(u.workdir, fmt.Sprintf("%s.pid", SERVICE)),
	}

	directories := []string{
		u.logdir,
		u.workdir,
	}

	for _, f := range files {
		fmt.Printf("   ... removing '%s'\n", f)
		if err := os.Remove(f); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	for _, dir := range directories {
		fmt.Printf("   ... removing '%s'\n", dir)
		if err := os.Remove(dir); err != nil && !os.IsNotExist(err) {
			patherr, ok := err.(*os.PathError)
			if !ok {
				return err
			}

			syserr, ok := patherr.Err.(syscall.Errno)
			if !ok {
				return err
			}

			if syserr != syscall.ENOTEMPTY {
				return err
			}

			fmt.Printf("   ... WARNING: could not remove directory '%s' (%v)\n", dir, syserr)
		}
	}

	return nil
}

func (u *Undaemonize) firewall(executable string) error {
	if executable == "" {
		return nil
	}

	fmt.Println()
	fmt.Println("   ***")
	fmt.Printf("   *** WARNING: removing '%s' from the application firewall\n", SERVICE)
	fmt.Println("   ***")
	fmt.Println()

	path := executable
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

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--remove", path)
		out, err = cmd.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("Failed to remove 'uhppoted-rest' from the application firewall (%v)\n", err)
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
