package commands

import (
	"flag"
	"os"
)

type Context struct {
}

type Command interface {
	Parse([]string) error
	Execute(context Context) error
	Cmd() string
	Description() string
	Usage() string
	Help()
}

var VERSION = "v0.04.0"

var cli = []Command{
	NewDaemonize(),
	NewUndaemonize(),
	&Version{VERSION},
	&Help{},
}

func Parse() (Command, error) {
	if len(os.Args) > 1 {
		for _, c := range cli {
			if c.Cmd() == flag.Arg(0) {
				return c, c.Parse(flag.Args()[1:])
			}
		}
	}

	return nil, nil
}
