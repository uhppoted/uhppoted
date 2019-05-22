package simulator

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net"
	"time"
	"uhppote"
	"uhppote-simulator/simulator/entities"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

type Simulator struct {
	File         string              `json:"-"`
	Compressed   bool                `json:"-"`
	SerialNumber types.SerialNumber  `json:"serial-number"`
	IpAddress    net.IP              `json:"address"`
	SubnetMask   net.IP              `json:"subnet"`
	Gateway      net.IP              `json:"gateway"`
	MacAddress   entities.MacAddress `json:"MAC"`
	Version      types.Version       `json:"version"`
	Date         types.Date          `json:"-"`
	Cards        entities.CardList   `json:"cards"`
}

type CardNotFoundResponse struct {
	MsgType      types.MsgType      `uhppote:"value:0x5a"`
	SerialNumber types.SerialNumber `uhppote:"offset:4"`
	CardNumber   uint32             `uhppote:"offset:8"`
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

func (s *Simulator) PutCard(request uhppote.PutCardRequest) (*uhppote.PutCardResponse, error) {
	card := entities.Card{
		CardNumber: request.CardNumber,
		From:       request.From,
		To:         request.To,
		Door1:      request.Door1,
		Door2:      request.Door2,
		Door3:      request.Door3,
		Door4:      request.Door4,
	}

	s.Cards.Put(&card)

	saved := false
	err := s.Save()
	if err == nil {
		saved = true
	}

	response := uhppote.PutCardResponse{
		SerialNumber: s.SerialNumber,
		Succeeded:    saved,
	}

	return &response, nil
}

func (s *Simulator) GetCards(request uhppote.GetCardsRequest) (*uhppote.GetCardsResponse, error) {
	response := uhppote.GetCardsResponse{
		SerialNumber: s.SerialNumber,
		Records:      uint32(len(s.Cards)),
	}

	return &response, nil
}

func (s *Simulator) GetCardById(request uhppote.GetCardByIdRequest) (interface{}, error) {
	for _, card := range s.Cards {
		if request.CardNumber == card.CardNumber {
			response := uhppote.GetCardByIdResponse{
				SerialNumber: s.SerialNumber,
				CardNumber:   card.CardNumber,
				From:         card.From,
				To:           card.To,
				Door1:        card.Door1,
				Door2:        card.Door2,
				Door3:        card.Door3,
				Door4:        card.Door4,
			}

			return &response, nil
		}
	}

	return &CardNotFoundResponse{
		SerialNumber: s.SerialNumber,
		CardNumber:   0,
	}, nil
}
