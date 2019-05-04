package commands

import (
	"fmt"
	"uhppote"
)

type ListenCommand struct {
}

func (c *ListenCommand) Execute(u *uhppote.UHPPOTE) error {
	fmt.Printf("Listening...\n")

	err := u.Listen()
	if err == nil {
	}

	return err
}

func (c *ListenCommand) CLI() string {
	return "listen"
}

func (c *ListenCommand) Description() string {
	return "Listens for access control events"
}

func (c *ListenCommand) Usage() string {
	return "listen"
}

func (c *ListenCommand) Help() {
	fmt.Println("Listens for access control events from UHPPOTE UTC3110-L0x controllers configured to send events to this IP address and port")
	fmt.Println()
}
