package commands

import (
	"flag"
	"fmt"
	"os"
)

type Context struct {
}

type Command interface {
	Name() string
	FlagSet() *flag.FlagSet
	Execute(context Context) error
	Description() string
	Usage() string
	Help()
}

var VERSION = "v0.4.2"

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
			if os.Args[1] == c.Name() {
				cmd = c
				args = os.Args[2:]
				break
			}
		}
	}

	flagset := cmd.FlagSet()
	if flagset == nil {
		panic(fmt.Sprintf("'%s' command implementation without a flagset: %#v", cmd.Name, cmd))
	}

	return cmd, flagset.Parse(args)
}
