package commands

import (
	"errors"
	"flag"
	"fmt"
	"time"
	"uhppote"
	"uhppote/types"
)

type SetTimeCommand struct {
	SerialNumber uint32
	DateTime     types.DateTime
}

func NewSetTimeCommand() (*SetTimeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	datetime := time.Now()
	if len(flag.Args()) > 2 {
		if flag.Arg(2) == "now" {
			datetime = time.Now()
		} else {
			datetime, err = time.Parse("2006-01-02 15:04:05", flag.Arg(2))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid date/time parameter: %v", flag.Arg(3)))
			}
		}
	}

	return &SetTimeCommand{serialNumber, types.DateTime{datetime}}, nil
}

func (c *SetTimeCommand) Execute(u *uhppote.UHPPOTE) error {
	devicetime, err := u.SetTime(c.SerialNumber, c.DateTime)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}
