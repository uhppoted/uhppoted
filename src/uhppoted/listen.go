package uhppoted

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"uhppote"
	"uhppote/types"
)

type EventMap struct {
	filepath  string
	retrieved map[uint32]uint32
}

type ListenEvent struct {
	DeviceID   uint32         `json:"device-id"`
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

func (u *UHPPOTED) Listen(ctx context.Context, received *EventMap, q chan os.Signal) {
	for device, index := range received.retrieved {
		event, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, 0xffffffff)
		if err != nil {
			u.warn(0, "listen", err)
		} else {
			if retrieved := u.fetch(ctx, device, index+1, event.Index); retrieved != 0 {
				received.retrieved[device] = retrieved
				if err := received.store(); err != nil {
					u.warn(0, "listen", err)
				}
			}
		}
	}

	u.listen(ctx, received, q)
}

func (u *UHPPOTED) listen(ctx context.Context, received *EventMap, q chan os.Signal) {
	p := make(chan *types.Status)

	go func() {
		if err := ctx.Value("uhppote").(*uhppote.UHPPOTE).Listen(p, q); err != nil {
			u.warn(0, "listen", err)
		}
	}()

	for {
		event := <-p
		if event == nil {
			break
		}

		u.log("EVENT", uint32(event.SerialNumber), fmt.Sprintf("%v", event))

		device := uint32(event.SerialNumber)
		last := event.LastIndex
		first := last
		retrieved, ok := received.retrieved[device]
		if ok {
			first = retrieved + 1
		}

		if retrieved := u.fetch(ctx, device, first, last); retrieved != 0 {
			received.retrieved[device] = retrieved
			if err := received.store(); err != nil {
				u.warn(0, "listen", err)
			}
		}
	}
}

func (u *UHPPOTED) fetch(ctx context.Context, device uint32, first uint32, last uint32) uint32 {
	retrieved := uint32(0)

	for index := first; index <= last; index++ {
		record, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetEvent(device, index)
		if err != nil {
			u.warn(device, "listen", fmt.Errorf("Failed to retrieve event ID %d", index))
			continue
		}

		if record == nil {
			u.warn(device, "listen", fmt.Errorf("No event record for ID %d", index))
			continue
		}

		if record.Index != index {
			u.warn(device, "listen", fmt.Errorf("No event record for ID %d", index))
			continue
		}

		retrieved = record.Index
		message := EventMessage{
			Event: ListenEvent{
				DeviceID:   device,
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
		u.send(ctx, message)
	}

	return retrieved
}

func NewEventMap(filepath string) *EventMap {
	return &EventMap{
		filepath:  filepath,
		retrieved: map[uint32]uint32{},
	}
}

func (m *EventMap) Load(log *log.Logger) error {
	if m.filepath == "" {
		return nil
	}

	f, err := os.Open(m.filepath)
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
	if m.filepath == "" {
		return nil
	}

	dir := path.Dir(m.filepath)
	filename := path.Base(m.filepath) + ".tmp"
	tmpfile := path.Join(dir, filename)

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

	return os.Rename(tmpfile, m.filepath)
}
