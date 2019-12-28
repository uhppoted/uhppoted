package commands

import (
	"fmt"
)

type GetCardCommand struct {
}

func (c *GetCardCommand) Execute(ctx Context) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return err
	}

	record, err := ctx.uhppote.GetCardByID(serialNumber, cardNumber)
	if err != nil {
		return err
	}

	if record == nil {
		fmt.Printf("%v %v NO RECORD\n", serialNumber, cardNumber)
	} else {
		fmt.Printf("%-10d %v\n", serialNumber, record)
	}

	return nil
}

func (c *GetCardCommand) CLI() string {
	return "get-card"
}

func (c *GetCardCommand) Description() string {
	return "Returns the access granted to a card number"
}

func (c *GetCardCommand) Usage() string {
	return "<serial number> <card number>"
}

func (c *GetCardCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-card <serial number> <card number>")
	fmt.Println()
	fmt.Println(" Retrieves the access granted for the card number from  the controller card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  card-number    (required) card number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-card 12345678 9876543")
	fmt.Println()
}
