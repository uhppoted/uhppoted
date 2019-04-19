package commands

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"uhppote"
)

type GrantCommand struct {
}

func (c *GrantCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	from, err := getDate(3, "Missing start date", "Invalid start date: %v")
	if err != nil {
		return err
	}

	to, err := getDate(4, "Missing end date", "Invalid end date: %v")
	if err != nil {
		return err
	}

	permissions, err := getPermissions(5)
	if err != nil {
		return err
	}

	authorised, err := u.Authorise(serialNumber, cardNumber, *from, *to, *permissions)

	if err == nil {
		fmt.Printf("%v\n", authorised)
	}

	return err
}

func (c *GrantCommand) CLI() string {
	return "grant"
}

func (c *GrantCommand) Description() string {
	return "Grants access to a card"
}

func (c *GrantCommand) Usage() string {
	return "<serial number> <card number> <start date> <end date> <doors>"
}

func (c *GrantCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] authorise <serial number> <card number> <start date> <end date> <doors>")
	fmt.Println()
	fmt.Println(" Adds a card to the authorised list")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println("  <start date>     (required) start date YYYY-MM-DD")
	fmt.Println("  <end date>       (required) end date   YYYY-MM-DD")
	fmt.Println("  <doors>          (required) list of permitted doors [1 2 3 4]")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli authorise 12345678 918273645 2019-01-01 2019-12-31 1,2,4")
	fmt.Println()
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
