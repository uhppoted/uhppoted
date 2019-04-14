package types

import "fmt"

type SwipeIndex struct {
	SerialNumber SerialNumber
	Index        uint32
}

type SwipeIndexResult struct {
	SerialNumber SerialNumber
	Index        uint32
	Succeeded    bool
}

type Swipe struct {
	SerialNumber SerialNumber
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
	return fmt.Sprintf("%s %d", s.SerialNumber, s.Index)
}

func (s *SwipeIndexResult) String() string {
	return fmt.Sprintf("%s %-8d %v", s.SerialNumber, s.Index, s.Succeeded)
}

func (s *Swipe) String() string {
	return fmt.Sprintf("%s %-4d %s %-12d %1d %-5v %-4d", s.SerialNumber, s.Index, s.Timestamp.String(), s.CardNumber, s.Door, s.Granted, s.RecordType)
}
