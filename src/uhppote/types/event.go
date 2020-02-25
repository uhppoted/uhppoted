package types

import "fmt"

type EventIndex struct {
	SerialNumber SerialNumber
	Index        uint32
}

type EventIndexResult struct {
	SerialNumber SerialNumber
	Index        uint32
	Changed      bool
}

type Event struct {
	SerialNumber SerialNumber
	Index        uint32
	Type         byte
	Granted      bool
	Door         byte
	DoorOpened   bool
	UserID       uint32
	Timestamp    DateTime
	Result       byte
}

func (s *EventIndex) String() string {
	return fmt.Sprintf("%s %d", s.SerialNumber, s.Index)
}

func (s *EventIndexResult) String() string {
	return fmt.Sprintf("%s %-8d %v", s.SerialNumber, s.Index, s.Changed)
}

func (s *Event) String() string {
	return fmt.Sprintf("%s %-4d %s %-12d %1d %-5v %-4d", s.SerialNumber, s.Index, s.Timestamp.String(), s.UserID, s.Door, s.Granted, s.Result)
}
