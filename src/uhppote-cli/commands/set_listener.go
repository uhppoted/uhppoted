package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"uhppote"
)

type SetListenerCommand struct {
}

func (c *SetListenerCommand) Execute(ctx context.Context, u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	if len(flag.Args()) < 3 {
		return errors.New("Missing IP address")
	}

	address, err := net.ResolveUDPAddr("udp", flag.Arg(2))
	if err != nil {
		return err
	}

	if address == nil || address.IP.To4() == nil {
		return errors.New(fmt.Sprintf("Invalid UDP address: %v", flag.Arg(2)))
	}

	listener, err := u.SetListener(serialNumber, *address)

	if err == nil {
		fmt.Printf("%v\n", listener)
	}

	return err
}

func (c *SetListenerCommand) CLI() string {
	return "set-listener"
}

func (c *SetListenerCommand) Description() string {
	return "Sets the IP address and port to which the controller sends access events"
}

func (c *SetListenerCommand) Usage() string {
	return "<serial number> <address:port>"
}

func (c *SetListenerCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] set-listener <serial number> <address:port>")
	fmt.Println()
	fmt.Println(" Sets the host address to which the controller sends access events")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  address        (required) IPv4 address")
	fmt.Println("  port           (required) IP port in the range 1 to 65535")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli set-listener 12345678  192.168.1.100:54321")
	fmt.Println()
}
