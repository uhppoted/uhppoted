package parsers

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"uhppote-cli/config"
	"uhppote/types"
)

type ACL struct {
	Path   string
	Config *config.Config
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
}

func (a *ACL) Parse(f *bufio.Reader) (*map[uint32][]types.Card, error) {
	acl := make(map[uint32][]types.Card)
	for id, _ := range a.Config.Devices {
		acl[id] = make([]types.Card, 0)
	}

	r := csv.NewReader(f)
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	index, err := parseHeader(header, a.Path, a.Config)
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

		cards, err := parseRecord(record, index, a.Path)
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

func parseRecord(record []string, index *index, path string) (*map[uint32]types.Card, error) {
	cards := make(map[uint32]types.Card, 0)
	for k, v := range index.doors {
		card := types.Card{Doors: make([]bool, 4)}

		if cardnumber, err := strconv.ParseUint(record[index.cardnumber-1], 10, 32); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid card number: '%s'", record[index.cardnumber-1]))
		} else {
			card.CardNumber = uint32(cardnumber)
		}

		if date, err := time.ParseInLocation("2006-01-02", record[index.from-1], time.Local); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid 'from' date: '%s'", record[index.from-1]))
		} else {
			card.From = types.Date(date)
		}

		if date, err := time.ParseInLocation("2006-01-02", record[index.to-1], time.Local); err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid 'to' date: '%s'", record[index.to-1]))
		} else {
			card.To = types.Date(date)
		}

		for i, d := range v {
			if d == 0 {
				card.Doors[i] = false
			} else {
				switch record[d-1] {
				case "Y":
					card.Doors[i] = true
				case "N":
					card.Doors[i] = false
				default:
					return nil, errors.New(fmt.Sprintf("Expected 'Y/N' for door: '%s'", record[d]))
				}
			}
		}

		cards[k] = card
	}

	return &cards, nil
}
