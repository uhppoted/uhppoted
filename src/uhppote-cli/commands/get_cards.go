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

func (c *GetCardsCommand) CLI() string {
	return "get-cards"
}

func (c *GetCardsCommand) Description() string {
	return "Returns the list of cards stored on the controller"
}

func (c *GetCardsCommand) Usage() string {
	return "<serial number>"
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
