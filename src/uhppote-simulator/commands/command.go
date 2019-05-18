package commands

type Command interface {
	Execute(dir string) error
	CLI() string
	Description() string
	Usage() string
	Help()
}
