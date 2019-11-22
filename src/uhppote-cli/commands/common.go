package commands

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getUint8(index int, missing, invalid string) (uint8, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 8)

	if err != nil {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return uint8(N), err
}

func getUint16(index int, missing, invalid string) (uint16, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 16)

	if err != nil {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return uint16(N), err
}

func getUint32(index int, missing, invalid string) (uint32, error) {
	if len(flag.Args()) < index+1 {
		return 0, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]+", flag.Arg(index))

	if !valid {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	N, err := strconv.ParseUint(flag.Arg(index), 10, 32)

	if err != nil {
		return 0, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return uint32(N), err
}

func getString(index int, missing, invalid string) (string, error) {
	if len(flag.Args()) < index+1 {
		return "", errors.New(missing)
	}

	return flag.Arg(index), nil
}

func getDate(index int, missing, invalid string) (*time.Time, error) {
	if len(flag.Args()) < index+1 {
		return nil, errors.New(missing)
	}

	valid, _ := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", flag.Arg(index))

	if !valid {
		return nil, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	date, err := time.Parse("2006-01-02", flag.Arg(index))

	if err != nil {
		return nil, errors.New(fmt.Sprintf(invalid, flag.Arg(index)))
	}

	return &date, err
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

func getPermissions(index int) ([]bool, error) {
	doors := []bool{false, false, false, false}

	if len(flag.Args()) > index {
		matches := strings.Split(flag.Arg(index), ",")

		for _, match := range matches {
			door, err := strconv.Atoi(match)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Invalid door '%v'", match))
			}

			if door > 0 && door < 5 {
				doors[door-1] = true
			}

		}
	}

	return doors, nil
}
