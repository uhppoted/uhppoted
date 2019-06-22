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

type Simulator struct {
	File         string             `json:"-"`
	Compressed   bool               `json:"-"`
	SerialNumber types.SerialNumber `json:"serial-number"`
	IpAddress    net.IP             `json:"address"`
	SubnetMask   net.IP             `json:"subnet"`
	Gateway      net.IP             `json:"gateway"`
	MacAddress   types.MacAddress   `json:"MAC"`
	Version      types.Version      `json:"version"`
	Date         types.Date         `json:"-"`
	Cards        entities.CardList  `json:"cards"`
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

	date, err := time.ParseInLocation("20060102", "20180816", time.Local)
	if err != nil {
		return nil, err
	}

	simulator.File = filepath
	simulator.Compressed = true
	simulator.Date = types.Date(date)

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

	date, err := time.ParseInLocation("20060102", "20180816", time.Local)
	if err != nil {
		return nil, err
	}

	simulator.File = filepath
	simulator.Compressed = false
	simulator.Date = types.Date(date)

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
