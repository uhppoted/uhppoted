package entities

import (
	"uhppote/types"
)

type Event struct {
	RecordNumber uint32         `json:"record-number"`
	Type         uint8          `json:"type"`
	Granted      bool           `json:"granted"`
	Door         uint8          `json:"door"`
	DoorOpened   bool           `json:"door-opened"`
	UserId       uint32         `json:"user-id"`
	Timestamp    types.DateTime `json:"timestamp"`
	RecordType   uint8          `json:"record-type"`
}

type EventList struct {
	Index  uint32  `json:"index"`
	Events []Event `json:"events"`
}

// TODO: implement Marshal/Unmarshal
func (l *EventList) Add(event *Event) {
	if event != nil {
		event.RecordNumber = uint32(len(l.Events) + 1)
		l.Events = append(l.Events, *event)
	}
}

func (l *EventList) Get(index uint32) *Event {
	if index > 0 && int(index) <= len(l.Events) {
		return &l.Events[index-1]
	}

	return nil
}

func (l *EventList) LastIndex() uint32 {
	return uint32(len(l.Events))
}

func (l *EventList) SetIndex(index uint32) bool {
	if index != l.Index {
		l.Index = index
		return true
	}

	return false
}
