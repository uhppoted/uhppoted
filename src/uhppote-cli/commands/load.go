package commands

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"uhppote-cli/config"
	"uhppote-cli/parsers"
	"uhppote/types"
)

type Load struct {
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
}

type mergelist struct {
	id     uint32
	ignore []types.Card
	add    []types.Card
	update []types.Card
	delete []types.Card
}

func (c *Load) Execute(ctx Context) error {
	err := ctx.config.Verify()
	if err != nil {
		return err
	}

	file, err := getACLFile()
	if err != nil {
		return err
	}

	acl, err := parse(*file, ctx.config)
	if err != nil {
		return err
	}

	list := make([]mergelist, 0)

	for id, cards := range *acl {
		device, err := getCards(ctx, id)
		if err != nil {
			return err
		}

		list = append(list, diff(id, cards, *device))
	}

	for _, rs := range list {
		fmt.Println("  ID:     ", rs.id)
		fmt.Println("  IGNORE: ", rs.ignore)
		fmt.Println("  ADD:    ", rs.add)
		fmt.Println("  UPDATE: ", rs.update)
		fmt.Println("  DELETE: ", rs.delete)
	}

	return nil
}

func getACLFile() (*string, error) {
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

func parse(path string, cfg *config.Config) (*parsers.ACL, error) {
	fmt.Printf("   ... loading access control list from '%s'\n", path)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	acl := parsers.ACL{}

	return acl.Load(bufio.NewReader(f), path, cfg)
}

func getCards(ctx Context, serialNumber uint32) (*[]types.Card, error) {
	N, err := ctx.uhppote.GetCards(serialNumber)
	if err != nil {
		return nil, err
	}

	cards := make([]types.Card, 0)

	for index := uint32(0); index < N.Records; index++ {
		record, err := ctx.uhppote.GetCardByIndex(serialNumber, index+1)
		if err != nil {
			return nil, err
		}
		cards = append(cards, *record)
	}

	return &cards, nil
}

func diff(id uint32, master, device []types.Card) mergelist {
	p := make(map[uint32]*types.Card)
	q := make(map[uint32]*types.Card)

	for n, c := range master {
		p[c.CardNumber] = &master[n]
	}

	for n, c := range device {
		q[c.CardNumber] = &device[n]
	}

	m := mergelist{
		id:     id,
		ignore: make([]types.Card, 0),
		add:    make([]types.Card, 0),
		update: make([]types.Card, 0),
		delete: make([]types.Card, 0),
	}

	for n, c := range p {
		if q[n] == nil {
			m.add = append(m.add, *c)
		} else if reflect.DeepEqual(c, q[n]) {
			m.ignore = append(m.ignore, *c)
		} else {
			m.update = append(m.update, *c)
		}
	}

	for n, c := range q {
		if p[n] == nil {
			m.delete = append(m.delete, *c)
		}
	}

	return m
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
