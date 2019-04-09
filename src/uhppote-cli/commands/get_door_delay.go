package commands

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"uhppote"
)

type GetDoorDelayCommand struct {
}

func (c *GetDoorDelayCommand) Execute(u *uhppote.UHPPOTE) error {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return err
	}

	door, err := getDoor(2, "Missing door", "Invalid door: %v")
	if err != nil {
		return err
	}

	record, err := u.GetDoorDelay(serialNumber, door)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", record)

	return nil
}

func (c *GetDoorDelayCommand) CLI() string {
	return "get-door-delay"
}

func (c *GetDoorDelayCommand) Help() {
	fmt.Println("Usage: uhppote-cli [options] get-door-delay <serial number> <door>")
	fmt.Println()
	fmt.Println(" Retrieves the door open delay (in seconds)")
	fmt.Println()
	fmt.Println("  serial-number  (required) controller serial number")
	fmt.Println("  door           (required) door (1,2,3 or 4")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli get-door 12345678 3")
	fmt.Println()
}

func getDoor(index int, missing, invalid string) (byte, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[1-4]", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	door, err := strconv.Atoi(flag.Arg(index))

	if err != nil {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return byte(door), nil
}
