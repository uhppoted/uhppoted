package simulator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net"
	"time"
	"uhppote-simulator/simulator/entities"
	"uhppote/types"
)

type Offset time.Duration

type Door struct {
	Open   bool  `json:"-"`
	Button bool  `json:"-"`
	Delay  uint8 `json:"delay"`
}

type Simulator struct {
	File           string             `json:"-"`
	Compressed     bool               `json:"-"`
	SerialNumber   types.SerialNumber `json:"serial-number"`
	IpAddress      net.IP             `json:"address"`
	SubnetMask     net.IP             `json:"subnet"`
	Gateway        net.IP             `json:"gateway"`
	MacAddress     types.MacAddress   `json:"MAC"`
	Version        types.Version      `json:"version"`
	TimeOffset     Offset             `json:"offset"`
	Cards          entities.CardList  `json:"cards"`
	Events         entities.EventList `json:"events"`
	SystemState    byte               `json:"state"`
	Doors          map[uint8]Door     `json:"doors"`
	PacketNumber   uint32             `json:"packet-number"`
	Backup         uint32             `json:"backup"`
	SpecialMessage byte               `json:"special-message"`
	Battery        byte               `json:"battery"`
	FireAlarm      byte               `json:"fire-alarm"`
}

func (t Offset) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(t).String())
}

func (t *Offset) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*t = Offset(d)

	return nil
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
		simulator.Doors = make(map[uint8]Door)
		for i := 1; i <= 4; i++ {
			ix := uint8(i)
			simulator.Doors[ix] = Door{false, false, 5}
		}
	}

	return simulator, nil
}

func (s *Simulator) Save() error {
	if s.Compressed {
		return saveGZ(s.File, s)
	}

	return save(s.File, s)
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
