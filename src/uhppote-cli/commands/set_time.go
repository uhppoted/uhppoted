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
	Debug        bool
}

func NewSetTimeCommand(debug bool) (*SetTimeCommand, error) {
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

	return &SetTimeCommand{serialNumber, types.DateTime{datetime}, debug}, nil
}

func (c *SetTimeCommand) Execute() error {
	u := uhppote.UHPPOTE{SerialNumber: c.SerialNumber, Debug: c.Debug}
	devicetime, err := u.SetTime(c.DateTime)

	if err == nil {
		fmt.Printf("%s\n", devicetime)
	}

	return err
}
