package commands

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"os"
	"path/filepath"
	"syscall"
)

var UNDAEMONIZE = Undaemonize{
	name:    SERVICE,
	workdir: workdir(),
	logdir:  filepath.Join(workdir(), "logs"),
	config:  workdir(),
}

type Undaemonize struct {
	name    string
	workdir string
	logdir  string
	config  string
}

func (u *Undaemonize) Name() string {
	return "undaemonize"
}

func (u *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (u *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters %s from the list of Windows services", SERVICE)
}

func (u *Undaemonize) Usage() string {
	return ""
}

func (u *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s from the list of Windows services", SERVICE)
	fmt.Println()
}

func (u *Undaemonize) Execute(ctx context.Context) error {
	fmt.Println("   ... undaemonizing")

	if err := u.unregister(); err != nil {
		return err
	}

	if err := u.clean(); err != nil {
		return err
	}

	fmt.Printf("   ... %s deregistered as a Windows service\n", SERVICE)
	fmt.Printf(`
   NOTE: Configuration files in %s,
               working files in %s,
               and log files in %s
               were not removed and should be deleted manually
`, filepath.Dir(u.config), u.workdir, u.logdir)
	fmt.Println()

	return nil
}

func (u *Undaemonize) unregister() error {
	fmt.Println("   ... unregistering %s as a Windows service", u.name)
	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(u.name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", u.name)
	}

	defer s.Close()

	fmt.Printf("   ... stopping %s service\n", u.name)
	status, err := s.Control(svc.Stop)
	if err != nil {
		return err
	}
	fmt.Printf("   ... %s stopped: %v\n", u.name, status)

	fmt.Printf("   ... deleting %s service\n", u.name)
	err = s.Delete()
	if err != nil {
		return err
	}

	err = eventlog.Remove(u.name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}

	fmt.Printf("   ... %s unregistered from the list of Windows services\n", u.name)
	return nil
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

			// Windows error is: ERROR_DIR_NOT_EMPTY (0x91). May be fixed in 1.14.
			// Ref. https://github.com/golang/go/issues/32309
			// Ref. https://docs.microsoft.com/en-us/windows/win32/debug/system-error-codes--0-499-
			if syserr != syscall.ENOTEMPTY && syserr != 0x91 {
				return err
			}

			fmt.Printf("   ... WARNING: could not remove directory '%s' (%v)\n", dir, syserr)
		}
	}

	return nil
}
