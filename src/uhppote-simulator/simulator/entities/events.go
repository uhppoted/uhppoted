package entities

import (
	"uhppote/types"
)

type Event struct {
	RecordNumber byte           `json:"record-number"`
	Granted      bool           `json:"granted"`
	Door         byte           `json:"door"`
	Opened       bool           `json:"opened"`
	CardNumber   uint32         `json:"card-number"`
	Timestamp    types.DateTime `json:"timestamp"`
	Reason       byte           `json:"reason"`
	Door1State   bool           `json:"door1-state"`
	Door2State   bool           `json:"door2-state"`
	Door3State   bool           `json:"door3-state"`
	Door4State   bool           `json:"door4-state"`
	Door1Button  bool           `json:"door1-button"`
	Door2Button  bool           `json:"door2-button"`
	Door3Button  bool           `json:"door3-button"`
	Door4Button  bool           `json:"door4-button"`
}

type EventList struct {
	LastIndex uint32  `json:"index"`
	Events    []Event `json:"events"`
}

// TODO: implement Marshal/Unmarshal
func (l *EventList) Add(event *Event) {
	if event != nil {
		l.Events = append(l.Events, *event)
	}
}

func (l *EventList) Get(index uint32) *Event {
	if index > 0 && int(index) <= len(l.Events) {
		return &l.Events[index-1]
	}

	return nil
}
