package types

import (
	"fmt"
)

type Card struct {
	CardNumber uint32
	From       Date
	To         Date
	Doors      []bool
}

type Authorised struct {
	SerialNumber uint32
	Authorised   bool
}

func (c Card) String() string {
	f := func(d bool) string {
		if d {
			return "Y"
		}
		return "N"
	}

	return fmt.Sprintf("%-8v %v %v %s %s %s %s", c.CardNumber, c.From, c.To, f(c.Doors[0]), f(c.Doors[1]), f(c.Doors[2]), f(c.Doors[3]))
}

func (r *Authorised) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Authorised)
}
