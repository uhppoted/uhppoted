package commands

import (
	"fmt"
	"uhppote"
)

type ListDevicesCommand struct {
}

func NewListDevicesCommand() (*ListDevicesCommand, error) {
	return &ListDevicesCommand{}, nil
}

func (c *ListDevicesCommand) Execute(u *uhppote.UHPPOTE) error {
	devices, err := u.Search()

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}
