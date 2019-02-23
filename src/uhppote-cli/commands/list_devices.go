package commands

import (
	"fmt"
	"uhppote"
)

type ListDevicesCommand struct {
	Debug bool
}

func NewListDevicesCommand(debug bool) (*ListDevicesCommand, error) {
	return &ListDevicesCommand{debug}, nil
}

func (c *ListDevicesCommand) Execute() error {
	u := uhppote.UHPPOTE{Debug: c.Debug}
	devices, err := u.Search()

	if err == nil {
		for _, device := range devices {
			fmt.Printf("%s\n", device.String())
		}
	}

	return err
}
