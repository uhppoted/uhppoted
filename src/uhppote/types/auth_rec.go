package types

import "fmt"

type AuthRec struct {
	SerialNumber uint32
	Records      uint64
}

func (r *AuthRec) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Records)
}
