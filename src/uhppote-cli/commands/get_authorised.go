package commands

import (
	"fmt"
	"uhppote"
)

type GetAuthorisedCommand struct {
	SerialNumber uint32
	Debug        bool
}

func NewGetAuthorisedCommand(debug bool) (*GetAuthorisedCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetAuthorisedCommand{serialNumber, debug}, nil
}

func (c *GetAuthorisedCommand) Execute() error {
	u := uhppote.UHPPOTE{SerialNumber: c.SerialNumber, Debug: c.Debug}
	authorised, err := u.GetAuthRec()

	if err == nil {
		fmt.Printf("%v\n", authorised)
	}

	return err
}
