package uhppoted

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"uhppote/types"
)

type EventMap struct {
	file      string
	retrieved map[uint32]uint32
}

type ListenEvent struct {
	DeviceID   DeviceID       `json:"device-id"`
	EventID    uint32         `json:"event-id"`
	Type       uint8          `json:"event-type"`
	Granted    bool           `json:"access-granted"`
	Door       uint8          `json:"door-id"`
	DoorOpened bool           `json:"door-opened"`
	UserID     uint32         `json:"user-id"`
	Timestamp  types.DateTime `json:"timestamp"`
	Result     uint8          `json:"event-result"`
}

type EventMessage struct {
	Event ListenEvent `json:"event"`
}

type EventHandler func(EventMessage)

type listener struct {
	onConnected func()
	onEvent     func(*types.Status)
	onError     func(error) bool
}

func (l *listener) OnConnected() {
	go func() {
		l.onConnected()
	}()
}

func (l *listener) OnEvent(event *types.Status) {
	go func() {
		l.onEvent(event)
	}()
}

func (l *listener) OnError(err error) bool {
	return l.onError(err)
}

func (u *UHPPOTED) Listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	for device, index := range received.retrieved {
		event, err := u.Uhppote.GetEvent(device, 0xffffffff)
		if err != nil {
			u.warn("listen", err)
		} else {
			if retrieved := u.fetch(device, index+1, event.Index, handler); retrieved != 0 {
				received.retrieved[device] = retrieved
				if err := received.store(); err != nil {
					u.warn("listen", err)
				}
			}
		}
	}

	u.listen(handler, received, q)
}

func (u *UHPPOTED) listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	backoffs := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		5 * time.Second,
		10 * time.Second,
		20 * time.Second,
		30 * time.Second,
		60 * time.Second,
	}

	ix := 0
	l := listener{
		onConnected: func() {
			u.info("listen", "Connected")
			ix = 0
		},

		onEvent: func(e *types.Status) {
			u.onEvent(e, received, handler)
		},

		onError: func(err error) bool {
			u.warn("listen", err)
			return true
		},
	}

	for {
		if err := u.Uhppote.Listen(&l, q); err != nil {
			u.warn("listen", err)

			delay := 60 * time.Second
			if ix < len(backoffs) {
				delay = backoffs[ix]
				ix++
			}

			u.info("listen", fmt.Sprintf("Retrying in %v", delay))
			time.Sleep(delay)
		}
	}
}

func (u *UHPPOTED) onEvent(e *types.Status, received *EventMap, handler EventHandler) {
	u.info("event", fmt.Sprintf("%+v", e))

	device := uint32(e.SerialNumber)
	last := e.LastIndex
	first := last
	retrieved, ok := received.retrieved[device]
	if ok && retrieved < last {
		first = retrieved + 1
	}

	if eventID := u.fetch(device, first, last, handler); retrieved != 0 {
		received.retrieved[device] = eventID
		if err := received.store(); err != nil {
			u.warn("listen", err)
		}
	}
}

func (u *UHPPOTED) fetch(device uint32, first uint32, last uint32, handler EventHandler) uint32 {
	retrieved := uint32(0)

	for index := first; index <= last; index++ {
		record, err := u.Uhppote.GetEvent(device, index)
		if err != nil {
			u.warn("listen", fmt.Errorf("Failed to retrieve event for device %d, ID %d", device, index))
			continue
		}

		if record == nil {
			u.warn("listen", fmt.Errorf("No event record for device %d, ID %d", device, index))
			continue
		}

		if record.Index != index {
			u.warn("listen", fmt.Errorf("No event record for device %d, ID %d", device, index))
			continue
		}

		retrieved = record.Index
		message := EventMessage{
			Event: ListenEvent{
				DeviceID:   DeviceID(record.SerialNumber),
				EventID:    record.Index,
				Type:       record.Type,
				Granted:    record.Granted,
				Door:       record.Door,
				DoorOpened: record.DoorOpened,
				UserID:     record.UserID,
				Timestamp:  record.Timestamp,
				Result:     record.Result,
			},
		}

		u.debug("listen", fmt.Sprintf("event %v", message))
		handler(message)
	}

	return retrieved
}

func NewEventMap(file string) *EventMap {
	return &EventMap{
		file:      file,
		retrieved: map[uint32]uint32{},
	}
}

func (m *EventMap) Load(log *log.Logger) error {
	if m.file == "" {
		return nil
	}

	f, err := os.Open(m.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer f.Close()

	re := regexp.MustCompile(`^\s*(.*?)(?::\s*|\s*=\s*|\s+)(\S.*)\s*`)
	s := bufio.NewScanner(f)
	for s.Scan() {
		match := re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			if device, err := strconv.ParseUint(key, 10, 32); err != nil {
				log.Printf("WARN: Error parsing event map entry '%s': %v", s.Text(), err)
			} else if eventID, err := strconv.ParseUint(value, 10, 32); err != nil {
				log.Printf("WARN: Error parsing event map entry '%s': %v", s.Text(), err)
			} else {
				m.retrieved[uint32(device)] = uint32(eventID)
			}
		}
	}

	return s.Err()
}

func (m *EventMap) store() error {
	if m.file == "" {
		return nil
	}

	dir := filepath.Dir(m.file)
	filename := filepath.Base(m.file) + ".tmp"
	tmpfile := filepath.Join(dir, filename)

	f, err := os.Create(tmpfile)
	if err != nil {
		return err
	}

	defer f.Close()

	for key, value := range m.retrieved {
		if _, err := fmt.Fprintf(f, "%-16d %v\n", key, value); err != nil {
			return err
		}
	}

	f.Close()

	return os.Rename(tmpfile, m.file)
}
