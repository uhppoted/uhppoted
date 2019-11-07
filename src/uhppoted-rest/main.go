package main

import (
	"fmt"
	"os"
	"uhppoted-rest/commands"
)

func main() {
	cmd, err := commands.Parse()
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if cmd != nil {
		ctx := commands.Context{}
		if err = cmd.Execute(ctx); err != nil {
			fmt.Printf("\nERROR: %v\n\n", err)
			os.Exit(1)
		}
	}
}
