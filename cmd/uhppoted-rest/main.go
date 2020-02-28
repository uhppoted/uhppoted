package main

import (
	"context"
	"fmt"
	"github.com/uhppoted/uhppoted-api/command"
	"github.com/uhppoted/uhppoted-rest/commands"
	"os"
)

var cli = []uhppoted.Command{
	commands.NewDaemonize(),
	commands.NewUndaemonize(),
	&uhppoted.VERSION,
}

var help = uhppoted.NewHelp(commands.SERVICE, cli, &commands.RUN)

func main() {
	cmd, err := uhppoted.Parse(cli, &commands.RUN, help)
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err = cmd.Execute(ctx); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}
