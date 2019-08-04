package simulator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"reflect"
	"time"
	"uhppote"
	"uhppote-simulator/simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

type Simulator struct {
	File       string                `json:"-"`
	Compressed bool                  `json:"-"`
	TxQ        chan entities.Message `json:"-"`

	SerialNumber   types.SerialNumber       `json:"serial-number"`
	IpAddress      net.IP                   `json:"address"`
	SubnetMask     net.IP                   `json:"subnet"`
	Gateway        net.IP                   `json:"gateway"`
	MacAddress     types.MacAddress         `json:"MAC"`
	Version        types.Version            `json:"version"`
	TimeOffset     entities.Offset          `json:"offset"`
	Doors          map[uint8]*entities.Door `json:"doors"`
	Listener       *net.UDPAddr             `json:"listener"`
	SystemState    byte                     `json:"state"`
	PacketNumber   uint32                   `json:"packet-number"`
	Backup         uint32                   `json:"backup"`
	SpecialMessage byte                     `json:"special-message"`
	Battery        byte                     `json:"battery"`
	FireAlarm      byte                     `json:"fire-alarm"`
	Cards          entities.CardList        `json:"cards"`
	Events         entities.EventList       `json:"events"`
}

func (s *Simulator) Handle(src *net.UDPAddr, rq messages.Request) {
	s.send(src, s.handle(rq))
}

func (s *Simulator) handle(rq messages.Request) messages.Response {
	switch v := rq.(type) {
	case *messages.GetStatusRequest:
		return s.getStatus(rq.(*messages.GetStatusRequest))

	case *messages.SetTimeRequest:
		return s.setTime(rq.(*messages.SetTimeRequest))

	case *messages.GetTimeRequest:
		return s.getTime(rq.(*messages.GetTimeRequest))

	case *messages.OpenDoorRequest:
		return s.openDoor(rq.(*messages.OpenDoorRequest))

	case *messages.PutCardRequest:
		return s.putCard(rq.(*messages.PutCardRequest))

	case *messages.DeleteCardRequest:
		return s.deleteCard(rq.(*messages.DeleteCardRequest))

	case *messages.DeleteCardsRequest:
		return s.deleteCards(rq.(*messages.DeleteCardsRequest))

	case *messages.GetCardsRequest:
		return s.getCards(rq.(*messages.GetCardsRequest))

	case *messages.GetCardByIndexRequest:
		return s.getCardByIndex(rq.(*messages.GetCardByIndexRequest))

	case *messages.SetDoorDelayRequest:
		return s.setDoorDelay(rq.(*messages.SetDoorDelayRequest))

	case *messages.GetDoorDelayRequest:
		return s.GetDoorDelay(rq.(*messages.GetDoorDelayRequest))

	case *messages.SetListenerRequest:
		return s.setListener(rq.(*messages.SetListenerRequest))

	case *messages.GetListenerRequest:
		return s.getListener(rq.(*messages.GetListenerRequest))

	case *messages.FindDevicesRequest:
		return s.find(rq.(*messages.FindDevicesRequest))

	case *messages.SetAddressRequest:
		return s.setAddress(rq.(*messages.SetAddressRequest))

	case *messages.GetEventRequest:
		return s.getEvent(rq.(*messages.GetEventRequest))

	case *messages.SetEventIndexRequest:
		return s.setEventIndex(rq.(*messages.SetEventIndexRequest))

	case *messages.GetEventIndexRequest:
		return s.GetEventIndex(rq.(*messages.GetEventIndexRequest))

	default:
		panic(errors.New(fmt.Sprintf("Unsupported message type %T", v)))
	}

}

func Load(filepath string, compressed bool) (*Simulator, error) {
	if compressed {
		return loadGZ(filepath)
	}

	return load(filepath)
}

func loadGZ(filepath string) (*Simulator, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	buffer, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	simulator := new(Simulator)
	err = json.Unmarshal(buffer, simulator)
	if err != nil {
		return nil, err
	}

	simulator.File = filepath
	simulator.Compressed = true

	return simulator, nil
}

func load(filepath string) (*Simulator, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	simulator := new(Simulator)
	err = json.Unmarshal(bytes, simulator)
	if err != nil {
		return nil, err
	}

	simulator.File = filepath
	simulator.Compressed = false

	if simulator.Doors == nil {
		simulator.Doors = make(map[uint8]*entities.Door)
	}

	for i := uint8(1); i <= 4; i++ {
		if simulator.Doors[i] == nil {
			simulator.Doors[i] = entities.NewDoor(i)
		}
	}

	return simulator, nil
}

func (s *Simulator) Save() error {
	if s.File != "" {
		if s.Compressed {
			return saveGZ(s.File, s)
		}

		return save(s.File, s)
	}

	return nil
}

func (s *Simulator) send(dest *net.UDPAddr, message interface{}) {
	if dest != nil && message != nil && !reflect.ValueOf(message).IsNil() {
		s.TxQ <- entities.Message{dest, message}
	}
}

func saveGZ(filepath string, s *Simulator) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	zw := gzip.NewWriter(&buffer)
	_, err = zw.Write(b)
	if err != nil {
		return err
	}

	if err = zw.Close(); err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, buffer.Bytes(), 0644)
}

func save(filepath string, s *Simulator) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, bytes, 0644)
}

func (s *Simulator) add(e *entities.Event) {
	if e != nil {
		s.Events.Add(e)
		s.Save()

		utc := time.Now().UTC()
		datetime := utc.Add(time.Duration(s.TimeOffset))
		event := uhppote.Event{
			SerialNumber:   s.SerialNumber,
			LastIndex:      s.Events.LastIndex(),
			SystemState:    s.SystemState,
			SystemDate:     types.SystemDate(datetime),
			SystemTime:     types.SystemTime(datetime),
			PacketNumber:   s.PacketNumber,
			Backup:         s.Backup,
			SpecialMessage: s.SpecialMessage,
			LowBattery:     s.Battery,
			FireAlarm:      s.FireAlarm,

			Door1State: s.Doors[1].IsOpen(),
			Door2State: s.Doors[2].IsOpen(),
			Door3State: s.Doors[3].IsOpen(),
			Door4State: s.Doors[4].IsOpen(),

			Door1Button: s.Doors[1].IsButtonPressed(),
			Door2Button: s.Doors[2].IsButtonPressed(),
			Door3Button: s.Doors[3].IsButtonPressed(),
			Door4Button: s.Doors[4].IsButtonPressed(),

			// SwipeRecord =
			Granted:       e.Granted,
			Door:          e.Door,
			DoorOpened:    e.DoorOpened,
			UserId:        e.UserId,
			SwipeDateTime: e.Timestamp,
			SwipeReason:   e.Type,
		}

		s.send(s.Listener, event)
	}
}

func (s *Simulator) onError(err error) {
	fmt.Printf("ERROR: %v\n", err)
}
