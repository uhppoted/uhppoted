package commands

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"uhppote"
)

type SetAddressCommand struct {
	SerialNumber uint32
	Address      net.IP
	Mask         net.IP
	Gateway      net.IP
}

func NewSetAddressCommand() (*SetAddressCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	if len(flag.Args()) < 3 {
		return nil, errors.New("Missing IP address")
	}

	address := net.ParseIP(flag.Arg(2))

	if address == nil || address.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP address: %v", flag.Arg(2)))
	}

	mask := net.IPv4(255, 255, 255, 0)
	if len(flag.Args()) > 3 {
		mask = net.ParseIP(flag.Arg(3))

		if mask == nil || mask.To4() == nil {
			mask = net.IPv4(255, 255, 255, 0)
		}
	}

	gateway := net.IPv4(0, 0, 0, 0)

	if len(flag.Args()) > 4 {
		gateway = net.ParseIP(flag.Arg(3))
		if gateway == nil || gateway.To4() == nil {
			gateway = net.IPv4(0, 0, 0, 0)
		}
	}

	return &SetAddressCommand{serialNumber, address, mask, gateway}, nil
}

func (c *SetAddressCommand) Execute(u *uhppote.UHPPOTE) error {
	err := u.SetAddress(c.SerialNumber, c.Address, c.Mask, c.Gateway)

	return err
}

func (c *SetAddressCommand) CLI() string {
	return "set-address"
}

func (c *SetAddressCommand) Description() string {
	return "Sets the controller IP address"
}

func (c *SetAddressCommand) Usage() string {
	return "<serial number> <address> [mask] [gateway]"
}

func (c *SetAddressCommand) Help() {
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
