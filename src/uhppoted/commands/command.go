package commands

type Context struct {
}

type Command interface {
	Execute(context Context) error
	CLI() string
	Description() string
	Usage() string
	Help()
}
