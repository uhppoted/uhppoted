package types

import "fmt"

type SwipeCount struct {
	SerialNumber uint32
	Count        uint32
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

func (s *SwipeCount) String() string {
	return fmt.Sprintf("%-12d %d", s.SerialNumber, s.Count)
}

func (s *Swipe) String() string {
	return fmt.Sprintf("%-12d %-4d %s  %-8d %1d %-5v %d", s.SerialNumber, s.Index, s.Timestamp.String(), s.CardNumber, s.Door, s.Granted, s.RecordType)
}
