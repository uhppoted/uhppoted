package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
)

// DARWIN

var dir = flag.String("dir", "/usr/local/var/uhppoted", "Working directory")
var logfile = flag.String("logfile", "/usr/local/var/uhppoted/logs/uhppoted.log", "uhppoted log file")
var logfilesize = flag.Int("logfilesize", 10, "uhppoted log file size")
var pidFile = flag.String("pid", "/usr/local/var/uhppoted/uhppoted.pid", "uhppoted PID file")
var useSyslog = flag.Bool("syslog", false, "Use syslog for event logging")
var modifyFirewall = flag.Bool("modify-application-firewall", false, "Add 'uhppoted' to the MacOS application firewall and unblock incoming connections (requires 'sudo')")

func sysinit() {
	log.Printf("uhppoted daemon    - %s (PID %d)\n", "MacOS", os.Getpid())

	if *modifyFirewall {
		unblock()
	}
}

// NOTE: this is potentially fragile and an INTERIM workaround to run the uhppoted daemon on MacOS with the
//       MacOS application firewall enabled (as it should be!). Requires superuser (sudo) privileges so if
//       you use this option you're (hopefully) aware of the risks and are also (hopefully) running behind
//       an external firewall already. At the very least the REST endpoint should be HTTPS only with pinned
//       SSL client and server certificates.
func unblock() {
	log.Println()
	log.Println("   ***")
	log.Println("   *** WARNING: adding 'uhppoted' to the application firewall and unblocking incoming connections")
	log.Println("   ***")
	log.Println()

	path, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get path to executable: %v\n", err)
		return
	}

	cmd := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := cmd.CombinedOutput()
	log.Printf("   > %s", out)
	if err != nil {
		log.Fatalf("ERROR: Failed to retrieve application firewall global state (%v)\n", err)
		return
	}

	if strings.Contains(string(out), "State = 1") {
		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to disable the application firewall (%v)\n", err)
			return
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to add 'uhppoted' to the application firewall (%v)\n", err)
			return
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to unblock 'uhppoted' on the application firewall (%v)\n", err)
			return
		}

		cmd = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = cmd.CombinedOutput()
		log.Printf("   > %s", out)
		if err != nil {
			log.Fatalf("ERROR: Failed to re-enable the application firewall (%v)\n", err)
			return
		}

		log.Println()
	}
}
