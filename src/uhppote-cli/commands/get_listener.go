package commands

import (
	"context"
	"fmt"
	"uhppote"
)

type GetListenerCommand struct {
}

func (c *GetListenerCommand) Execute(ctx context.Context, u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	address, err := u.GetListener(serialNumber)

	if err == nil {
		fmt.Printf("%v\n", address)
	}

	return err
}

func (c *GetListenerCommand) CLI() string {
	return "get-listener"
}

func (c *GetListenerCommand) Description() string {
	return "Returns the IP address to which the selected controller sends events"
}

func (c *GetListenerCommand) Usage() string {
	return "<serial number>"
}

func (c *GetListenerCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-listener <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the IP address and port of the remote host to which the controller sends access events")
	fmt.Println(" with the corresponding serial number in the format:")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Example:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-listener 12345678")
	fmt.Println()
}
