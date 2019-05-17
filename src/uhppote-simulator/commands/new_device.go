package commands

import (
	"flag"
	"fmt"
	"net"
	"path"
	"uhppote-simulator/simulator"
	"uhppote/types"
)

type NewDeviceCommand struct {
}

func (c *NewDeviceCommand) Execute() error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	dir, err := getDir(2, "Missing devices directory", "Invalid devices directory: %v")
	if err != nil {
		return err
	}

	gzip := false
	filename := fmt.Sprintf("%d.json", serialNumber)
	if len(flag.Args()) > 3 {
		if flag.Arg(3) == "--gzip" {
			gzip = true
			filename = fmt.Sprintf("%d.json.gz", serialNumber)
		}
	}

	mac, _ := net.ParseMAC("00:66:19:39:55:2d")

	device := simulator.Simulator{
		SerialNumber: types.SerialNumber(serialNumber),
		IpAddress:    net.IPv4(192, 168, 0, 25),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(0, 0, 0, 0),
		MacAddress:   simulator.MacAddress(mac),
		Version:      0x0892,
	}

	if gzip {
		return simulator.SaveGZ(path.Join(dir, filename), &device)
	}

	return simulator.Save(path.Join(dir, filename), &device)
}

func (c *NewDeviceCommand) CLI() string {
	return "new-device"
}

func (c *NewDeviceCommand) Description() string {
	return "Creates a new simulator device file"
}

func (c *NewDeviceCommand) Usage() string {
	return "<serial number> <directory>"
}

func (c *NewDeviceCommand) Help() {
	fmt.Println("Usage: uhppote-simulator [options] new-device <serial number> <directory>")
	fmt.Println()
	fmt.Println(" Creates a new simulator device file in the specified directory")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  directory      (required) directory for simulator device files")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-simulator new-device 12345678 /etc/uhppote/simulator/devices")
	fmt.Println()
}