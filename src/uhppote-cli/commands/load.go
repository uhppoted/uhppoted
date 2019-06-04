package commands

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"uhppote-cli/config"
)

type Load struct {
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
}

type card struct {
	cardnumber uint32
	from       time.Time
	to         time.Time
	doors      []bool
}

func (c *Load) Execute(ctx Context) error {
	err := ctx.config.Verify()
	if err != nil {
		return err
	}

	file, err := getTSVFile()
	if err != nil {
		return err
	}

	acl, err := parse(*file, ctx.config)
	if err != nil {
		return err
	}

	for id, cards := range *acl {
		fmt.Println(id, cards)
		err = getCards(ctx, id)
		if err != nil {
			return err
		}
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

func parse(path string, cfg *config.Config) (*map[uint32][]card, error) {
	fmt.Printf("   ... loading access control list from '%s'\n", path)

	acl := make(map[uint32][]card)
	for id, _ := range cfg.Devices {
		acl[id] = make([]card, 0)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	index, err := parseHeader(header, path, cfg)
	if err != nil {
		return nil, err
	}
	line := 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		line += 1

		cards, err := parseRecord(record, index, path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Line %d: %v\n", line, err))
		}

		for id, card := range *cards {
			if acl[id] != nil {
				acl[id] = append(acl[id], card)
			}
		}
	}

	return &acl, nil
}

func parseHeader(header []string, path string, cfg *config.Config) (*index, error) {
	columns := make(map[string]int)

	index := index{
		cardnumber: 0,
		from:       0,
		to:         0,
		doors:      make(map[uint32][]int),
	}

	for id, _ := range cfg.Devices {
		index.doors[id] = make([]int, 4)
	}

	for c, field := range header {
		key := strings.ReplaceAll(strings.ToLower(field), " ", "")
		ix := c + 1

		if columns[key] != 0 {
			return nil, errors.New(fmt.Sprintf("Duplicate column name '%s' in File '%s", field, path))
		}

		columns[key] = ix
	}

	index.cardnumber = columns["cardnumber"]
	index.from = columns["from"]
	index.to = columns["to"]

	for id, device := range cfg.Devices {
		for i, door := range device.Door {
			if d := strings.ReplaceAll(strings.ToLower(door), " ", ""); d != "" {
				index.doors[id][i] = columns[d]
			}
		}
	}

	if index.cardnumber == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'Card Number'", path))
	}

	if index.from == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'From'", path))
	}

	if index.to == 0 {
		return nil, errors.New(fmt.Sprintf("File '%s' does not include a column 'to'", path))
	}

	for id, device := range cfg.Devices {
		for i, door := range device.Door {
			if d := strings.ReplaceAll(strings.ToLower(door), " ", ""); d != "" {
				if index.doors[id][i] == 0 {
					return nil, errors.New(fmt.Sprintf("File '%s' does not include a column for door '%s'", path, door))
				}
			}
		}
	}

	return &index, nil
}

func parseRecord(record []string, index *index, path string) (*map[uint32]card, error) {
	cards := make(map[uint32]card, 0)
	for k, v := range index.doors {
		card := card{doors: make([]bool, 4)}

		if cardnumber, err := strconv.ParseUint(record[index.cardnumber-1], 10, 32); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid card number: '%s'", record[index.cardnumber-1]))
		} else {
			card.cardnumber = uint32(cardnumber)
		}

		if date, err := time.ParseInLocation("2006-01-02", record[index.from-1], time.Local); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid 'from' date: '%s'", record[index.from-1]))
		} else {
			card.from = date
		}

		if date, err := time.ParseInLocation("2006-01-02", record[index.to-1], time.Local); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid 'to' date: '%s'", record[index.to-1]))
		} else {
			card.to = date
		}

		for i, d := range v {
			if d == 0 {
				card.doors[i] = false
			} else {
				switch record[d] {
				case "Y":
					card.doors[i] = true
				case "N":
					card.doors[i] = false
				default:
					return nil, errors.New(fmt.Sprintf("Expected 'Y/N' for door: '%s'", record[d]))
				}
			}
		}

		cards[k] = card
	}

	return &cards, nil
}

func getCards(ctx Context, serialNumber uint32) error {
	N, err := ctx.uhppote.GetCards(serialNumber)
	if err != nil {
		return err
	}

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index+1)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", record)
	}

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
