package commands

import "uhppote"

type Command interface {
	Execute(u *uhppote.UHPPOTE) error
	CLI() string
	Description() string
	Usage() string
	Help()
}
