package commands

import (
	"fmt"
	"uhppote"
)

type GetTimeCommand struct {
	SerialNumber uint32
}

func NewGetTimeCommand() (*GetTimeCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	return &GetTimeCommand{serialNumber}, nil
}

func (c *GetTimeCommand) Execute(u *uhppote.UHPPOTE) error {
	datetime, err := u.GetTime(c.SerialNumber)

	if err == nil {
		fmt.Printf("%s\n", datetime)
	}

	return err
}

func (c *GetTimeCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-time <serial number> [command options]")
	fmt.Println()
	fmt.Println(" Retrieves the current date/time referenced to the local timezone for the access control board")
	fmt.Println(" with the corresponding serial number in the format:")
	fmt.Println()
	fmt.Println(" <serial number> <yyyy-mm-dd HH:mm:ss>")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    -debug  Displays vaguely useful internal information")
	fmt.Println()
	fmt.Println("  Command options:")
	fmt.Println()
}
