package UT0311L04

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"reflect"
	"time"
	"uhppote"
	"uhppote-simulator/entities"
	"uhppote/messages"
	"uhppote/types"
)

type UT0311L04 struct {
	file       string                `json:"-"`
	compressed bool                  `json:"-"`
	txq        chan entities.Message `json:"-"`

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

func NewUT0311L04(deviceId uint32, dir string, compressed bool) *UT0311L04 {
	filename := fmt.Sprintf("%d.json", deviceId)
	if compressed {
		filename = fmt.Sprintf("%d.json.gz", deviceId)
	}

	mac := make([]byte, 6)
	rand.Read(mac)

	device := UT0311L04{
		file:       path.Join(dir, filename),
		compressed: compressed,

		SerialNumber: types.SerialNumber(deviceId),
		IpAddress:    net.IPv4(0, 0, 0, 0),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(0, 0, 0, 0),
		MacAddress:   types.MacAddress(mac),
		Version:      0x0892,
		Doors: map[uint8]*entities.Door{
			1: entities.NewDoor(1),
			2: entities.NewDoor(2),
			3: entities.NewDoor(3),
			4: entities.NewDoor(4),
		},
	}

	return &device
}

func (s *UT0311L04) DeviceID() uint32 {
	return uint32(s.SerialNumber)
}

func (s *UT0311L04) DeviceType() string {
	return "UT0311-L04"
}

func (s *UT0311L04) FilePath() string {
	return s.file
}

func (s *UT0311L04) SetTxQ(txq chan entities.Message) {
	s.txq = txq
}

func (s *UT0311L04) Handle(src *net.UDPAddr, rq messages.Request) {
	switch v := rq.(type) {
	case *messages.GetStatusRequest:
		s.getStatus(src, rq.(*messages.GetStatusRequest))

	case *messages.SetTimeRequest:
		s.setTime(src, rq.(*messages.SetTimeRequest))

	case *messages.GetTimeRequest:
		s.getTime(src, rq.(*messages.GetTimeRequest))

	case *messages.OpenDoorRequest:
		s.openDoor(src, rq.(*messages.OpenDoorRequest))

	case *messages.PutCardRequest:
		s.putCard(src, rq.(*messages.PutCardRequest))

	case *messages.DeleteCardRequest:
		s.deleteCard(src, rq.(*messages.DeleteCardRequest))

	case *messages.DeleteCardsRequest:
		s.deleteCards(src, rq.(*messages.DeleteCardsRequest))

	case *messages.GetCardsRequest:
		s.getCards(src, rq.(*messages.GetCardsRequest))

	case *messages.GetCardByIdRequest:
		s.getCardById(src, rq.(*messages.GetCardByIdRequest))

	case *messages.GetCardByIndexRequest:
		s.getCardByIndex(src, rq.(*messages.GetCardByIndexRequest))

	case *messages.SetDoorControlStateRequest:
		s.setDoorDelay(src, rq.(*messages.SetDoorControlStateRequest))

	case *messages.GetDoorControlStateRequest:
		s.getDoorDelay(src, rq.(*messages.GetDoorControlStateRequest))

	case *messages.SetListenerRequest:
		s.setListener(src, rq.(*messages.SetListenerRequest))

	case *messages.GetListenerRequest:
		s.getListener(src, rq.(*messages.GetListenerRequest))

	case *messages.FindDevicesRequest:
		s.find(src, rq.(*messages.FindDevicesRequest))

	case *messages.SetAddressRequest:
		s.setAddress(src, rq.(*messages.SetAddressRequest))

	case *messages.GetEventRequest:
		s.getEvent(src, rq.(*messages.GetEventRequest))

	case *messages.SetEventIndexRequest:
		s.setEventIndex(src, rq.(*messages.SetEventIndexRequest))

	case *messages.GetEventIndexRequest:
		s.getEventIndex(src, rq.(*messages.GetEventIndexRequest))

	default:
		panic(errors.New(fmt.Sprintf("Unsupported message type %T", v)))
	}
}

func Load(filepath string, compressed bool) (*UT0311L04, error) {
	if compressed {
		return loadGZ(filepath)
	}

	return load(filepath)
}

func loadGZ(filepath string) (*UT0311L04, error) {
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

	simulator := new(UT0311L04)
	err = json.Unmarshal(buffer, simulator)
	if err != nil {
		return nil, err
	}

	simulator.file = filepath
	simulator.compressed = true

	return simulator, nil
}

func load(filepath string) (*UT0311L04, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	simulator := new(UT0311L04)
	err = json.Unmarshal(bytes, simulator)
	if err != nil {
		return nil, err
	}

	simulator.file = filepath
	simulator.compressed = false

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

func (s *UT0311L04) Save() error {
	if s.file != "" {
		if s.compressed {
			return saveGZ(s.file, s)
		}

		return save(s.file, s)
	}

	return nil
}

func (s *UT0311L04) Delete() error {
	if s.file != "" {
		if err := os.Remove(s.file); err != nil {
			return err
		}

		if _, err := os.Stat(s.file); err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

func (s *UT0311L04) send(dest *net.UDPAddr, message interface{}) {
	if s.txq == nil {
		panic(fmt.Sprintf("Device %d: missing TXQ", s.SerialNumber))
	}

	if s.txq != nil && dest != nil && message != nil && !reflect.ValueOf(message).IsNil() {
		s.txq <- entities.Message{dest, message}
	}
}

func saveGZ(filepath string, s *UT0311L04) error {
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

func save(filepath string, s *UT0311L04) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, bytes, 0644)
}

func (s *UT0311L04) add(e *entities.Event) uint32 {
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

			EventType:   e.Type,
			EventResult: e.Result,
			Timestamp:   e.Timestamp,
			UserId:      e.UserId,
			Granted:     e.Granted,
			Door:        e.Door,
			DoorOpened:  e.DoorOpened,
		}

		s.send(s.Listener, &event)

		return s.Events.LastIndex()
	}

	return 0
}
