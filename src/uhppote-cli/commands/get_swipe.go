package commands

import (
	"fmt"
	"uhppote"
)

type GetSwipeCommand struct {
	SerialNumber uint32
	Index        uint32
	Debug        bool
}

func NewGetSwipeCommand(debug bool) (*GetSwipeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	index, err := getUint32(2, "Missing swipe index", "Invalid swipe index: %v")
	if err != nil {
		return nil, err
	}

	return &GetSwipeCommand{serialNumber, index, debug}, nil
}

func (c *GetSwipeCommand) Execute() error {
	u := uhppote.UHPPOTE{SerialNumber: c.SerialNumber, Debug: c.Debug}
	swipe, err := u.GetSwipe(c.Index)

	if err == nil {
		if swipe != nil {
			fmt.Printf("%12d %s\n", u.SerialNumber, swipe.String())
		}
	}

	return err
}
