package commands

import "uhppote"

type Command interface {
	Execute(u *uhppote.UHPPOTE) error
	Help()
}
