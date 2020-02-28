package entities

import (
	"github.com/uhppoted/uhppote-core/types"
)

type Event struct {
	RecordNumber uint32         `json:"record-number"`
	Type         uint8          `json:"type"`
	Granted      bool           `json:"granted"`
	Door         uint8          `json:"door"`
	DoorOpened   bool           `json:"door-opened"`
	UserID       uint32         `json:"user-id"`
	Timestamp    types.DateTime `json:"timestamp"`
	Result       uint8          `json:"result"`
}

type EventList struct {
	Size   uint32  `json:"size"`
	First  uint32  `json:"first"`
	Last   uint32  `json:"last"`
	Index  uint32  `json:"index"`
	Events []Event `json:"events"`
}

// TODO: implement Marshal/Unmarshal
func (l *EventList) Add(event *Event) {
	if event != nil {
		l.Last = l.Last + 1
		if l.Last > l.Size {
			l.Last = 1
		}

		if l.Last == l.First {
			l.First = l.First + 1
			if l.First > l.Size {
				l.First = 1
			}
		}

		event.RecordNumber = uint32(l.Last)

		index := l.Last
		if index >= uint32(len(l.Events)) {
			l.Events = append(l.Events, *event)
		} else {
			l.Events[index-1] = *event
		}
	}
}

func (l *EventList) Get(index uint32) *Event {
	if len(l.Events) > 0 {
		if index == 0 {
			return &l.Events[l.First-1]
		}

		if index == 0xffffffff || index > uint32(len(l.Events)) {
			return &l.Events[l.Last-1]
		}

		if index > 0 && int(index) <= len(l.Events) {
			return &l.Events[index-1]
		}
	}

	return nil
}

func (l *EventList) SetIndex(index uint32) bool {
	if index != l.Index {
		l.Index = index
		return true
	}

	return false
}
