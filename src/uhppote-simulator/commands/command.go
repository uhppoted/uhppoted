package commands

type Command interface {
	Execute() error
	CLI() string
	Description() string
	Usage() string
	Help()
}
