package commands

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"uhppote-cli/config"
)

type Load struct {
}

func (c *Load) Execute(ctx Context) error {
	fmt.Printf("--------------- DEBUG:%v\n", ctx.uhppote)
	fmt.Printf("--------------- DEBUG:%v\n", ctx.config)

	file, err := getTSVFile()
	if err != nil {
		return err
	}

	err = parse(*file, ctx.config)
	if err != nil {
		return err
	}

	return nil
}

func getTSVFile() (*string, error) {
	if len(flag.Args()) < 2 {
		return nil, errors.New("ERROR: Please specify the TSV file from which to load the access control list ")
	}

	file := flag.Arg(1)
	stat, err := os.Stat(file)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(fmt.Sprintf("File '%s' does not exist", file))
		}

		return nil, errors.New(fmt.Sprintf("Failed to find file '%s':%v", file, err))
	}

	if stat.Mode().IsDir() {
		return nil, errors.New(fmt.Sprintf("File '%s' is a directory", file))
	}

	if !stat.Mode().IsRegular() {
		return nil, errors.New(fmt.Sprintf("File '%s' is not a real file", file))
	}

	return &file, nil
}

func parse(path string, cfg *config.Config) error {
	fmt.Printf("   ... loading access control list from '%s'\n", path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return err
	}

	columns := make(map[string]int)

	for c, field := range header {
		key := strings.ReplaceAll(strings.ToLower(field), " ", "")
		index := c + 1
		columns[key] = index
	}

	if columns["cardnumber"] == 0 {
		return errors.New(fmt.Sprintf("File '%s' does not include a column 'Card Number'", path))
	}

	if columns["from"] == 0 {
		return errors.New(fmt.Sprintf("File '%s' does not include a column 'From'", path))
	}

	if columns["to"] == 0 {
		return errors.New(fmt.Sprintf("File '%s' does not include a column 'to'", path))
	}

	fmt.Println(columns)

	doors := make(map[uint32][]int)
	for id, d := range cfg.Devices {
		doors[id] = make([]int, 4)

		if col := columns[strings.ReplaceAll(strings.ToLower(d.Door1), " ", "")]; col == 0 {
			return errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, d.Door1))
		} else {
			doors[id][0] = col
		}

		if col := columns[strings.ReplaceAll(strings.ToLower(d.Door2), " ", "")]; col == 0 {
			return errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, d.Door2))
		} else {
			doors[id][1] = col
		}

		if col := columns[strings.ReplaceAll(strings.ToLower(d.Door3), " ", "")]; col == 0 {
			return errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, d.Door3))
		} else {
			doors[id][2] = col
		}

		if col := columns[strings.ReplaceAll(strings.ToLower(d.Door4), " ", "")]; col == 0 {
			return errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, d.Door4))
		} else {
			doors[id][3] = col
		}
	}

	fmt.Println(doors)

	return nil
}

func (c *Load) CLI() string {
	return "load"
}

func (c *Load) Description() string {
	return "Downloads an access control list from a TSV file to a set of access controllers"
}

func (c *Load) Usage() string {
	return "<TSV file>"
}

func (c *Load) Help() {
	fmt.Println("Usage: uhppote-cli [options] load <TSV file>")
	fmt.Println()
	fmt.Println(" Downloads the access control list in the TSV file to the access controllers defined in the configuration file")
	fmt.Println()
	fmt.Println("  <TSV file>  (required) TSV file with access control list")
	fmt.Println("              The TSV file should conform to the following format:")
	fmt.Println("              Card Number<tab>From<tab>To<tab>Front Door<tab>Back Door<tab> ...")
	fmt.Println("              123456789<tab>2019-01-01<tab>2019-12-31<tab>Y<tab>N<tab> ...")
	fmt.Println("              987654321<tab>2019-03-05<tab>2019-11-15<tab>N<tab>N<tab> ...")
	fmt.Println()
	fmt.Println("              'Front Door', 'Back Door', etc should match the door labels in the configuration file.")
	fmt.Println("              The CLI will load the access control permissions across all the controllers listed,")
	fmt.Println("              adding cards where necessary and deleting cards not listed in the TSV file. Making")
	fmt.Println("              a backup copy of the existing permissions (using e.g. get-cards) before executing this")
	fmt.Println("              is highly recommended.")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println()
	fmt.Println("    uhppote-cli --config .config load \"hell-2019-05-25.tsv\"")
	fmt.Println()
}
