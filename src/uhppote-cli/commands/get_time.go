package commands

import (
	"fmt"
	"uhppote"
)

type GetTimeCommand struct {
	SerialNumber uint32
	Debug        bool
}

func NewGetTimeCommand(debug bool) (*GetTimeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetTimeCommand{serialNumber, debug}, nil
}

func (c *GetTimeCommand) Execute() error {
	u := uhppote.UHPPOTE{SerialNumber: c.SerialNumber, Debug: c.Debug}
	datetime, err := u.GetTime()

	if err == nil {
		fmt.Printf("%s\n", datetime)
	}

	return err
}
