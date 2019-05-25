package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"uhppote"
	"uhppote/types"
)

type ListenCommand struct {
}

func (c *ListenCommand) Execute(ctx context.Context, u *uhppote.UHPPOTE) error {
	fmt.Printf("Listening...\n")

	p := make(chan *types.Status)
	q := make(chan os.Signal)

	go func() {
		for {
			event := <-p
			fmt.Printf("%v\n", event)
		}
	}()

	signal.Notify(q, os.Interrupt)

	return u.Listen(p, q)
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
