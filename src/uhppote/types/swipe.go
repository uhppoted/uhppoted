package types

import "fmt"

type Swipe struct {
	Index      uint32
	Type       byte
	Access     bool
	Door       byte
	DoorState  byte
	Card       uint32
	Timestamp  DateTime
	RecordType byte
}

func (s *Swipe) String() string {
	return fmt.Sprintf("%4d: %s  %-8d %1d %-5v", s.Index, s.Timestamp.String(), s.Card, s.Door, s.Access)
}
