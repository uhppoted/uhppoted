package commands

import (
	"fmt"
	"uhppote"
)

type GetCardsCommand struct {
	SerialNumber uint32
}

func NewGetCardsCommand() (*GetCardsCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetCardsCommand{serialNumber}, nil
}

func (c *GetCardsCommand) Execute(u *uhppote.UHPPOTE) error {
	N, err := u.GetCardCount(c.SerialNumber)

	if err == nil {
		fmt.Printf("%v\n", N)
	}

	return err
}

func (c *GetCardsCommand) CLI() string {
	return "get-cards"
}

func (c *GetCardsCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-cards <serial number>")
	fmt.Println()
	fmt.Println(" Retrieves the number of cards in the controller card list")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-cards 12345678")
	fmt.Println()
}
