package types

import "fmt"

type SerialNumber uint32

type RecordCount struct {
	SerialNumber SerialNumber
	Records      uint32
}

func (s SerialNumber) String() string {
	return fmt.Sprintf("%-10d", s)
}

func (r *RecordCount) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Records)
}
