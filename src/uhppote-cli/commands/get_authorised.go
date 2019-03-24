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

func (c *GetAuthorisedCommand) CLI() string {
	return "get-authorised"
}

func (c *GetAuthorisedCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-authorised <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the number of authorised cards")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-authorised 12345678")
	fmt.Println()
}
