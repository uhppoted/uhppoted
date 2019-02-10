package types

import "fmt"

type Authorised struct {
	SerialNumber uint32
	Authorised   bool
}

func (r *Authorised) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Authorised)
}
