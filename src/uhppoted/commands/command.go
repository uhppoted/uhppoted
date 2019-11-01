package commands

import (
	"flag"
	"fmt"
	"os"
)

type Context struct {
}

type Command interface {
	FlagSet() *flag.FlagSet
	Parse([]string) error
	Execute(context Context) error
	Description() string
	Usage() string
	Help()
}

var VERSION = "v0.04.0"

var cli = []Command{
	NewDaemonize(),
	NewUndaemonize(),
	&version,
	&Help{},
}

// Can't invoke flag.Parse() because the 'run' is the default command and only want to parse the flags
// if the other commands are not being invoked
func Parse() (Command, error) {
	var cmd Command = &runCmd
	var args []string = os.Args[1:]

	if len(os.Args) > 1 {
		for _, c := range cli {
			flagset := c.FlagSet()
			if flagset == nil {
				panic(fmt.Sprintf("command without a flagset: %#v", c))
			}

			if flagset.Name() == os.Args[1] {
				cmd = c
				args = os.Args[1:]
				break
			}
		}
	}

	return cmd, cmd.Parse(args)
}
