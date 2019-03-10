package commands

import (
	"fmt"
	"uhppote"
)

type GetStatusCommand struct {
	SerialNumber uint32
}

func NewGetStatusCommand() (*GetStatusCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetStatusCommand{serialNumber}, nil
}

func (c *GetStatusCommand) Execute(u *uhppote.UHPPOTE) error {
	status, err := u.GetStatus(c.SerialNumber)

	if err == nil {
		fmt.Printf("%v\n", status)
	}

	return err
}
