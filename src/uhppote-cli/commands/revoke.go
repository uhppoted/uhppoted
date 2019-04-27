package commands

import (
	"fmt"
	"uhppote"
)

type RevokeCommand struct {
}

func (c *RevokeCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	result, err := u.DeleteCard(serialNumber, cardNumber)

	if err == nil {
		fmt.Printf("%v\n", result)
	}

	return err
}

func (c *RevokeCommand) CLI() string {
	return "revoke"
}

func (c *RevokeCommand) Description() string {
	return "Revokes access for a card"
}

func (c *RevokeCommand) Usage() string {
	return "<serial number> <card number>"
}

func (c *RevokeCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] revoke <serial number> <card number>")
	fmt.Println()
	fmt.Println(" Removes a card from the authorised list")
	fmt.Println()
	fmt.Println("  <serial number>  (required) controller serial number")
	fmt.Println("  <card number>    (required) card number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli revoke 12345678 918273645")
	fmt.Println()
}
