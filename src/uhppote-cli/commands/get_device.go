package commands

import (
	"fmt"
)

type GetDeviceCommand struct {
}

func (c *GetDeviceCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	device, err := ctx.uhppote.FindDevice(serialNumber)
	if err != nil {
		return err
	} else if device == nil {
		return fmt.Errorf("No device found matching serial number '%d'", serialNumber)
	}

	fmt.Printf("%s\n", device.String())

	return nil
}

func (c *GetDeviceCommand) CLI() string {
	return "get-device"
}

func (c *GetDeviceCommand) Description() string {
	return "'pings' a UHPPOTE controller using the IP address configured for the device"
}

func (c *GetDeviceCommand) Usage() string {
	return "<serial number>"
}

func (c *GetDeviceCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-device <serial number>>")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println(" Issues a 'get-devices' request directed at the IP address configured for the supplied serial number")
	fmt.Println(" and extracts the response matching the serial number, returning the board summary information in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <IP address> <subnet mask> <gateway> <MAC address> <hexadecimal version> <firmware date>")
	fmt.Println()
	fmt.Println(" Falls back to a broadcast on the local network if no IP address is configured for the supplied serial number.")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-device 12345678")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays a trace of request/response messages")
	fmt.Println()
}
