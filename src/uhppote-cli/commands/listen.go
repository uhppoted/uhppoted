package commands

import (
	"fmt"
	"os"
	"os/signal"
	"uhppote/types"
)

type ListenCommand struct {
}

func (c *ListenCommand) Execute(ctx Context) error {
	p := make(chan *types.Status)
	q := make(chan os.Signal)

	defer close(q)

	go func() {
		for {
			event := <-p
			fmt.Printf("%v\n", event)
		}
	}()

	signal.Notify(q, os.Interrupt)

	connected := func() {
		fmt.Printf("Listening...\n")
	}

	return ctx.uhppote.Listen(p, q, connected)
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
	fmt.Println("Listens for access control events from UHPPOTE UT0311-L0x controllers configured to send events to this IP address and port")
	fmt.Println()
}
