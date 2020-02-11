package commands

import (
	"fmt"
	"os"
	"strings"
	"uhppoted/config"
)

func dump(path string) error {
	fmt.Println()
	fmt.Printf("   ... displaying configuration information from '%s'\n", path)

	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	var s strings.Builder

	if err := cfg.Write(&s); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s\n", s.String())
	fmt.Println()

	return nil
}
