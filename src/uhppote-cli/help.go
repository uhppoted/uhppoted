package main

import "fmt"

func helpCommands() {
	fmt.Println("Supported commands:")
	fmt.Println()
	fmt.Println(" search")
	fmt.Println(" get-time <serial number>")
	fmt.Println()
}

func helpVersion() {
	fmt.Println("Displays the uhppote-cli version in the format v<major>.<minor>.<build> e.g. v1.00.10")
	fmt.Println()
}

func helpSearch() {
	fmt.Println("Usage: uhppote-cli [options] search [command options]")
	fmt.Println()
	fmt.Println(" Searches the local network for UHPPOTE access control boards reponding to a poll")
	fmt.Println(" on the default UDP port 60000. Returns a list of boards one per line in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <IP address> <subnet mask> <gateway> <MAC address> <hexadecimal version> <firmware date>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

func helpGetTime() {
	fmt.Println("Usage: uhppote-cli [options] get-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Retrieves the current date/time referenced to the local timezone for the access control board")
	fmt.Println(" with the corresponding serial number in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <yyyy-mm-dd HH:mm:ss>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}

func helpSetTime() {
	fmt.Println("Usage: uhppote-cli [options] set-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Sets the controller date/time to the supplied time. Defaults to 'now'. Command format")
	fmt.Println()
	fmt.Println(" <serial number> [now|<yyyy-mm-dd HH:mm:ss>]")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
	fmt.Println("    now                    Sets the controller time to the system time of the local system")
	fmt.Println("    'yyyy-mm-dd HH:mm:ss'  Sets the controller time to the explicitly supplied instant")
	fmt.Println()
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-time")
	fmt.Println("    uhppote-cli set-time now")
	fmt.Println("    uhppote-cli set-time '2019-01-12 20:15:32'")
	fmt.Println()
}

func helpSetAddress() {
	fmt.Println("Usage: uhppote-cli [options] set-address <serial number> <address> [mask] [gateway]")
	fmt.Println()
	fmt.Println(" Sets the controller IP address, subnet mask and gateway address")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  address        (required) IPv4 address")
	fmt.Println("  mask           (optional) IPv4 subnet mask. Defaults to 255.255.255.0")
	fmt.Println("  gateway        (optional) IPv4 gateway address. Defaults to 0.0.0.0")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-address 12345678  192.168.1.100")
	fmt.Println("    uhppote-cli set-address 12345678  192.168.1.100  255.255.255.0")
	fmt.Println("    uhppote-cli set-address 12345678  192.168.1.100  255.255.255.0  0.0.0.0")
	fmt.Println()
}

func helpGetAuthRec() {
	fmt.Println("Usage: uhppote-cli [options] get-auth-rec <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the authorised card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-autrh-rec 12345678")
	fmt.Println()
}
