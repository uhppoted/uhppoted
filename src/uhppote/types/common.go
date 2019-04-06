package types

import "fmt"

type RecordCount struct {
	SerialNumber uint32
	Records      uint32
}

func (r *RecordCount) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Records)
}
