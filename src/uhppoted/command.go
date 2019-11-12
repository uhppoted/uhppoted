package uhppoted

import (
	"context"
	"flag"
	"fmt"
	"os"
)

type Command interface {
	Name() string
	FlagSet() *flag.FlagSet
	Execute(context.Context) error
	Description() string
	Usage() string
	Help()
}

// Can't invoke flag.Parse() because the 'run' is the default command i.e. only want to parse the flags
// if the other commands are not being invoked
func Parse(cli []Command, run Command, help Command) (Command, error) {
	var cmd Command = run
	var args []string = os.Args[1:]

	if len(os.Args) > 1 {
		if os.Args[1] == help.Name() {
			cmd = help
			args = os.Args[2:]
		} else {
			for _, c := range cli {
				if os.Args[1] == c.Name() {
					cmd = c
					args = os.Args[2:]
					break
				}
			}
		}
	}

	flagset := cmd.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", cmd.Name, cmd))
	}

	return cmd, flagset.Parse(args)
}
