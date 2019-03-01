package commands

import (
	"fmt"
	"uhppote"
)

type GetTimeCommand struct {
	SerialNumber uint32
}

func NewGetTimeCommand() (*GetTimeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetTimeCommand{serialNumber}, nil
}

func (c *GetTimeCommand) Execute(u *uhppote.UHPPOTE) error {
	datetime, err := u.GetTime(c.SerialNumber)

	if err == nil {
		fmt.Printf("%s\n", datetime)
	}

	return err
}
