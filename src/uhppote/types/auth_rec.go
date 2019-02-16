package types

import "fmt"

type AuthRec struct {
	SerialNumber uint32
	Records      uint32
}

func (r *AuthRec) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Records)
}
