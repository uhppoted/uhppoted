package commands

import (
	"fmt"
)

type GetDevicesCommand struct {
}

func (c *GetDevicesCommand) Execute(ctx Context) error {
	devices, err := ctx.uhppote.FindDevices()

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}

func (c *GetDevicesCommand) CLI() string {
	return "get-devices"
}

func (c *GetDevicesCommand) Description() string {
	return "Returns a list of found UHPPOTE controllers on the network"
}

func (c *GetDevicesCommand) Usage() string {
	return ""
}

func (c *GetDevicesCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-devices [command options]")
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
}
