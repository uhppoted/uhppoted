package uhppoted

import (
	"errors"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"time"
)

const ROLLOVER = uint32(100000)

type GetEventRangeRequest struct {
	DeviceID DeviceID
	Start    *types.DateTime
	End      *types.DateTime
}

type GetEventRangeResponse struct {
	DeviceID DeviceID    `json:"device-id,omitempty"`
	Dates    *DateRange  `json:"dates,omitempty"`
	Events   *EventRange `json:"events,omitempty"`
}

type GetEventRequest struct {
	DeviceID DeviceID
	EventID  uint32
}

type GetEventResponse struct {
	DeviceID DeviceID `json:"device-id"`
	Event    event    `json:"event"`
}

type DateRange struct {
	Start *types.DateTime `json:"start,omitempty"`
	End   *types.DateTime `json:"end,omitempty"`
}

func (d *DateRange) String() string {
	if d.Start != nil && d.End != nil {
		return fmt.Sprintf("{ Start:%v, End:%v }", d.Start, d.End)
	}

	if d.Start != nil {
		return fmt.Sprintf("{ Start:%v }", d.Start)
	}

	if d.End != nil {
		return fmt.Sprintf("{ End:%v }", d.End)
	}

	return "{}"
}

type EventRange struct {
	First uint32 `json:"first"`
	Last  uint32 `json:"last"`
}

func (e *EventRange) String() string {
	return fmt.Sprintf("{ First:%v, Last:%v }", e.First, e.Last)
}

type EventIndex uint32

func (index EventIndex) increment(rollover uint32) EventIndex {
	ix := uint32(index)

	if ix < 1 {
		ix = 1
	} else if ix >= rollover {
		ix = 1
	} else {
		ix += 1
	}

	return EventIndex(ix)
}

func (index EventIndex) decrement(rollover uint32) EventIndex {
	ix := uint32(index)

	if ix <= 1 {
		ix = rollover
	} else if ix > rollover {
		ix = rollover
	} else {
		ix -= 1
	}

	return EventIndex(ix)
}

type event struct {
	Index      uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserID     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

func (u *UHPPOTED) GetEventRange(request GetEventRangeRequest) (*GetEventRangeResponse, error) {
	u.debug("get-events", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	start := request.Start
	end := request.End
	rollover := ROLLOVER

	if d, ok := u.Uhppote.Devices[device]; ok {
		if d.Rollover != 0 {
			rollover = d.Rollover
		}
	}

	f, err := u.Uhppote.GetEvent(device, 0)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting first event index from %v (%w)", device, err))
	} else if f == nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting first event index from %v (%w)", device, errors.New("Record not found")))
	}

	l, err := u.Uhppote.GetEvent(device, 0xffffffff)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting last event index from %v (%w)", device, err))
	} else if l == nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting last event index from %v (%w)", device, errors.New("Record not found")))
	}

	// The indexing logic below 'decrements' the index from l(ast) to f(irst) assuming that the on-device event store has
	// a circular event buffer of size ROLLOVER. The logic assumes the events are ordered by datetime, which is reasonable
	// but not necessarily true e.g. if the start/end interval includes a significant device time change.
	var first *types.Event
	var last *types.Event

	if start != nil || end != nil {
		index := EventIndex(l.Index)
		for {
			record, err := u.Uhppote.GetEvent(device, uint32(index))
			if err != nil {
				return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting event for index %v from %v (%w)", index, device, err))
			}

			if in(record, start, end) {
				if last == nil {
					last = record
				}

				first = record
			} else if first != nil || last != nil {
				break
			}

			if uint32(index) == f.Index {
				break
			}

			index = index.decrement(rollover)
		}
	}

	dates := (*DateRange)(nil)
	if start != nil || end != nil {
		dates = &DateRange{
			Start: start,
			End:   end,
		}
	}

	events := (*EventRange)(nil)
	if first != nil && last != nil {
		events = &EventRange{
			First: first.Index,
			Last:  last.Index,
		}
	}

	response := GetEventRangeResponse{
		DeviceID: DeviceID(device),
		Dates:    dates,
		Events:   events,
	}

	u.debug("get-events", fmt.Sprintf("response %+v", response))

	return &response, nil
}

func in(record *types.Event, start, end *types.DateTime) bool {
	if start != nil && time.Time(record.Timestamp).Before(time.Time(*start)) {
		return false
	}

	if end != nil && time.Time(record.Timestamp).After(time.Time(*end)) {
		return false
	}

	return true
}

func (u *UHPPOTED) GetEvent(request GetEventRequest) (*GetEventResponse, error) {
	u.debug("get-events", fmt.Sprintf("request  %+v", request))

	device := uint32(request.DeviceID)
	eventID := request.EventID

	record, err := u.Uhppote.GetEvent(device, eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", InternalServerError, fmt.Errorf("Error getting event for ID %v from %v (%w)", eventID, device, err))
	}

	if record == nil {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("No event record for ID %v for %v", eventID, device))
	}

	if record.Index != eventID {
		return nil, fmt.Errorf("%w: %v", NotFound, fmt.Errorf("No event record for ID %v for %v", eventID, device))
	}

	response := GetEventResponse{
		DeviceID: DeviceID(record.SerialNumber),
		Event: event{
			Index:      record.Index,
			Type:       record.Type,
			Granted:    record.Granted,
			Door:       record.Door,
			DoorOpened: record.DoorOpened,
			UserID:     record.UserID,
			Timestamp:  record.Timestamp,
			Result:     record.Result,
		},
	}

	u.debug("get-event", fmt.Sprintf("response %+v", response))

	return &response, nil
}
