package types

import "fmt"

type MsgType uint8

type Result struct {
	SerialNumber SerialNumber
	Succeeded    bool
}

type RecordCount struct {
	SerialNumber SerialNumber
	Records      uint32
}

func (r *Result) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Succeeded)
}

func (r *RecordCount) String() string {
	return fmt.Sprintf("%v %v", r.SerialNumber, r.Records)
}
