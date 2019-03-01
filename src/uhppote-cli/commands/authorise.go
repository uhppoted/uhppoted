package commands

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"uhppote"
	"uhppote/types"
)

type GrantCommand struct {
	SerialNumber uint32
	CardNumber   uint32
	From         types.Date
	To           types.Date
	Permissions  []int
}

func NewGrantCommand() (*GrantCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return nil, err
	}

	from, err := getDate(3, "Missing start date", "Invalid start date: %v")
	if err != nil {
		return nil, err
	}

	to, err := getDate(4, "Missing end date", "Invalid end date: %v")
	if err != nil {
		return nil, err
	}

	permissions, err := getPermissions(5)
	if err != nil {
		return nil, err
	}

	return &GrantCommand{serialNumber, cardNumber, *from, *to, *permissions}, nil
}

func (c *GrantCommand) Execute(u *uhppote.UHPPOTE) error {
	authorised, err := u.Authorise(c.SerialNumber, c.CardNumber, c.From, c.To, c.Permissions)

	if err == nil {
		fmt.Printf("%v\n", authorised)
	}

	return err
}

func getPermissions(index int) (*[]int, error) {
	permissions := []int{}

	if len(flag.Args()) > index {
		matches := strings.Split(flag.Arg(index), ",")

		for _, match := range matches {
			door, err := strconv.Atoi(match)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid door '%v'", match))
			}

			permissions = append(permissions, door)
		}
	}

	return &permissions, nil
}
