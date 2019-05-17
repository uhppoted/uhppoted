package simulator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"time"
	"uhppote"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

type MacAddress net.HardwareAddr
type Version types.Version

type Simulator struct {
	File         string             `json:"-"`
	SerialNumber types.SerialNumber `json:"serial-number"`
	IpAddress    net.IP             `json:"address"`
	SubnetMask   net.IP             `json:"subnet"`
	Gateway      net.IP             `json:"gateway"`
	MacAddress   MacAddress         `json:"MAC"`
	Version      Version            `json:"version"`
	Date         types.Date         `json:"-"`
}

func (m MacAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.HardwareAddr(m).String())
}

func (m *MacAddress) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	mac, err := net.ParseMAC(s)
	if err != nil {
		return err
	}

	*m = MacAddress(mac)

	return nil
}

func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("%04x", v))
}

func (v *Version) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	N, err := fmt.Sscanf(s, "%04x", v)
	if err != nil {
		return err
	}

	if N != 1 {
		return errors.New("Unable to extract 'version' from JSON file")
	}

	return nil
}

func LoadGZ(filepath string) (*Simulator, error) {
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
	simulator.Date = types.Date{date}

	return simulator, nil
}

func SaveGZ(filepath string, s *Simulator) error {
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

func Load(filepath string) (*Simulator, error) {
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
	simulator.Date = types.Date{date}

	return simulator, nil
}

func Save(filepath string, s *Simulator) error {
	bytes, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath, bytes, 0644)
}

func (s *Simulator) Find(bytes []byte) ([]byte, error) {
	response := uhppote.FindDevicesResponse{
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   net.HardwareAddr(s.MacAddress),
		Version:      types.Version(s.Version),
		Date:         s.Date,
	}

	reply, err := codec.Marshal(response)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *Simulator) GetCardById(bytes []byte) ([]byte, error) {
	from, _ := time.ParseInLocation("2006-01-02", "2019-02-03", time.Local)
	to, _ := time.ParseInLocation("2006-01-02", "2019-12-29", time.Local)

	response := uhppote.GetCardByIdResponse{
		SerialNumber: s.SerialNumber,
		CardNumber:   123456,
		From:         types.Date{from},
		To:           types.Date{to},
		Door1:        true,
		Door2:        false,
		Door3:        false,
		Door4:        true,
	}

	reply, err := codec.Marshal(response)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
