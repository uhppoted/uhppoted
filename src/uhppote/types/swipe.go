package types

import "fmt"

type SwipeIndex struct {
	SerialNumber uint32
	Index        uint32
}

type Swipe struct {
	SerialNumber uint32
	Index        uint32
	Type         byte
	Granted      bool
	Door         byte
	DoorState    byte
	CardNumber   uint32
	Timestamp    DateTime
	RecordType   byte
}

func (s *SwipeIndex) String() string {
	return fmt.Sprintf("%-12d %d", s.SerialNumber, s.Index)
}

func (s *Swipe) String() string {
	return fmt.Sprintf("%-12d %-4d %s %-12d %1d %-5v %-4d", s.SerialNumber, s.Index, s.Timestamp.String(), s.CardNumber, s.Door, s.Granted, s.RecordType)
}
