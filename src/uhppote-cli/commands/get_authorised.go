package commands

import (
	"fmt"
	"uhppote"
)

type GetAuthorisedCommand struct {
	SerialNumber uint32
}

func NewGetAuthorisedCommand() (*GetAuthorisedCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetAuthorisedCommand{serialNumber}, nil
}

func (c *GetAuthorisedCommand) Execute(u *uhppote.UHPPOTE) error {
	authorised, err := u.GetAuthRec(c.SerialNumber)

	if err == nil {
		fmt.Printf("%v\n", authorised)
	}

	return err
}
