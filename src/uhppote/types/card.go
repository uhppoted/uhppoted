package types

import (
	"fmt"
)

type Card struct {
	SerialNumber SerialNumber
	CardNumber   uint32
	From         Date
	To           Date
	Door1        bool
	Door2        bool
	Door3        bool
	Door4        bool
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

	return fmt.Sprintf("%s %-8v %v %v %s %s %s %s", c.SerialNumber, c.CardNumber, c.From, c.To, f(c.Door1), f(c.Door2), f(c.Door3), f(c.Door4))
}

func (r *Authorised) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Authorised)
}
