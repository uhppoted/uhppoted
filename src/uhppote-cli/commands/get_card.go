package commands

import (
	"fmt"
	"uhppote"
)

type GetCardCommand struct {
	SerialNumber uint32
	CardNumber   uint32
}

func NewGetCardCommand() (*GetCardCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	cardNumber, err := getUint32(2, "Missing card number", "Invalid card number: %v")
	if err != nil {
		return nil, err
	}

	return &GetCardCommand{serialNumber, cardNumber}, nil
}

func (c *GetCardCommand) Execute(u *uhppote.UHPPOTE) error {
	N, err := u.GetCards(c.SerialNumber)

	if err != nil {
		return err
	}

	for index := uint32(0); index < N.Records; index++ {
		record, err := u.GetCardByIndex(c.SerialNumber, index+1)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", record)
	}

	return nil
}

func (c *GetCardCommand) CLI() string {
	return "get-card"
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
