package uhppoted

import (
	"bufio"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
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

type EventHandler func(EventMessage) bool

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

const BATCHSIZE = 32

func (u *UHPPOTED) Listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	var wg sync.WaitGroup

	for _, d := range u.Uhppote.Devices {
		deviceID := d.DeviceID
		wg.Add(1)
		go func() {
			defer wg.Done()
			u.retrieve(deviceID, received, handler)
		}()
	}

	wg.Wait()

	u.listen(handler, received, q)
}

func (u *UHPPOTED) retrieve(deviceID uint32, received *EventMap, handler EventHandler) {
	if index, ok := received.retrieved[deviceID]; ok {
		u.info("listen", fmt.Sprintf("Fetching unretrieved events for device ID %v", deviceID))

		event, err := u.Uhppote.GetEvent(deviceID, 0xffffffff)
		if err != nil {
			u.warn("listen", fmt.Errorf("Unable to retrieve events for device ID %v (%w)", deviceID, err))
			return
		}

		if event.Index == uint32(index) {
			u.info("listen", fmt.Sprintf("No unretrieved events for device ID %v", deviceID))
			return
		}

		rollover := ROLLOVER
		if d, ok := u.Uhppote.Devices[deviceID]; ok {
			if d.Rollover != 0 {
				rollover = d.Rollover
			}
		}

		from := EventIndex(index)
		to := EventIndex(event.Index)

		if retrieved := u.fetch(deviceID, from.increment(rollover), to, handler); retrieved != 0 {
			received.retrieved[deviceID] = retrieved
			if err := received.store(); err != nil {
				u.warn("listen", err)
			}
		}
	}
}

func (u *UHPPOTED) listen(handler EventHandler, received *EventMap, q chan os.Signal) {
	u.info("listen", "Initialising event listener")

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
			u.info("listen", "Listening")
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

	// NTS: use 'for {..}' because 'for err := u.Uhppote.Listen; ..' only ever executes the
	//      'Listen' once - on loop initialization
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
			continue
		}

		break
	}
}

func (u *UHPPOTED) onEvent(e *types.Status, received *EventMap, handler EventHandler) {
	u.info("event", fmt.Sprintf("%+v", e))

	deviceID := uint32(e.SerialNumber)
	last := EventIndex(e.LastIndex)
	first := EventIndex(e.LastIndex)

	retrieved, ok := received.retrieved[deviceID]
	if ok && retrieved != uint32(last) {
		first = EventIndex(retrieved)
	}

	if eventID := u.fetch(deviceID, first, last, handler); eventID != 0 {
		received.retrieved[deviceID] = eventID
		if err := received.store(); err != nil {
			u.warn("listen", err)
		}
	}
}

func (u *UHPPOTED) fetch(deviceID uint32, from, to EventIndex, handler EventHandler) (retrieved uint32) {
	batchSize := BATCHSIZE
	rollover := ROLLOVER

	if u.ListenBatchSize > 0 {
		batchSize = u.ListenBatchSize
	}

	if d, ok := u.Uhppote.Devices[deviceID]; ok {
		if d.Rollover != 0 {
			rollover = d.Rollover
		}
	}

	first, err := u.Uhppote.GetEvent(deviceID, 0)
	if err != nil {
		u.warn("listen", fmt.Errorf("Failed to retrieve 'first' event for device %d (%w)", deviceID, err))
		return
	} else if first == nil {
		u.warn("listen", fmt.Errorf("No 'first' event record returned for device %d", deviceID))
		return
	}

	last, err := u.Uhppote.GetEvent(deviceID, 0xffffffff)
	if err != nil {
		u.warn("listen", fmt.Errorf("Failed to retrieve 'last' event for device %d (%w)", deviceID, err))
		return
	} else if first == nil {
		u.warn("listen", fmt.Errorf("No 'last' event record returned for device %d", deviceID))
		return
	}

	if last.Index >= first.Index {
		if uint32(from) < first.Index || uint32(from) > last.Index {
			from = EventIndex(first.Index)
		}

		if uint32(to) < first.Index || uint32(to) > last.Index {
			to = EventIndex(last.Index)
		}
	} else {
		if uint32(from) < first.Index && uint32(from) > last.Index {
			from = EventIndex(first.Index)
		}

		if uint32(to) < first.Index && uint32(to) > last.Index {
			to = EventIndex(last.Index)
		}
	}

	count := 0
	index := from
	for {
		count += 1
		if count > batchSize {
			return
		}

		record, err := u.Uhppote.GetEvent(deviceID, uint32(index))
		if err != nil {
			u.warn("listen", fmt.Errorf("Failed to retrieve event for device %d, ID %d (%w)", deviceID, index, err))
		} else if record == nil {
			u.warn("listen", fmt.Errorf("No event record for device %d, ID %d", deviceID, index))
		} else if record.Index != uint32(index) {
			u.warn("listen", fmt.Errorf("No event record for device %d, ID %d", deviceID, index))
		} else {
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
			if !handler(message) {
				break
			}

			retrieved = record.Index
		}

		if index == to {
			break
		}

		index = index.increment(rollover)
	}

	return
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
	if m.file == "" || IsDevNull(m.file) {
		return nil
	}

	f, err := ioutil.TempFile(os.TempDir(), "uhppoted*.tmp")
	if err != nil {
		return err
	}

	defer os.Remove(f.Name())

	for key, value := range m.retrieved {
		if _, err := fmt.Fprintf(f, "%-16d %v\n", key, value); err != nil {
			f.Close()
			return err
		}
	}

	f.Close()

	return os.Rename(f.Name(), m.file)
}
