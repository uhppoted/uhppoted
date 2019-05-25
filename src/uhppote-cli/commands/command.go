package commands

import (
	"context"
	"uhppote"
)

type Command interface {
	Execute(ctx context.Context, u *uhppote.UHPPOTE) error
	CLI() string
	Description() string
	Usage() string
	Help()
}
