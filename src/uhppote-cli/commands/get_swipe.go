package commands

import (
	"fmt"
	"uhppote"
)

type GetSwipeCommand struct {
	SerialNumber uint32
	Index        uint32
}

func NewGetSwipeCommand() (*GetSwipeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	index, err := getUint32(2, "Missing swipe index", "Invalid swipe index: %v")
	if err != nil {
		return nil, err
	}

	return &GetSwipeCommand{serialNumber, index}, nil
}

func (c *GetSwipeCommand) Execute(u *uhppote.UHPPOTE) error {
	swipe, err := u.GetSwipe(c.SerialNumber, c.Index)

	if err == nil {
		if swipe != nil {
			fmt.Printf("%12d %s\n", c.SerialNumber, swipe.String())
		}
	}

	return err
}
