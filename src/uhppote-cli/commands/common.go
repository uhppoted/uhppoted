package commands

import (
	"errors"
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

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
