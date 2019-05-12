package simulator

import (
	"net"
	"time"
	"uhppote"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

type Simulator struct {
	SerialNumber types.SerialNumber
	IpAddress    net.IP
	SubnetMask   net.IP
	Gateway      net.IP
	MacAddress   net.HardwareAddr
	Version      types.Version
	Date         types.Date
}

func NewSimulator(serialNo uint32) *Simulator {
	mac, _ := net.ParseMAC("00:66:19:39:55:2d")
	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)

	return &Simulator{
		SerialNumber: types.SerialNumber(serialNo),
		IpAddress:    net.IPv4(192, 168, 0, 25),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(0, 0, 0, 0),
		MacAddress:   mac,
		Version:      0x0892,
		Date:         types.Date{date},
	}
}

func (s *Simulator) Find(bytes []byte) ([]byte, error) {
	response := uhppote.FindDevicesResponse{
		SerialNumber: s.SerialNumber,
		IpAddress:    s.IpAddress,
		SubnetMask:   s.SubnetMask,
		Gateway:      s.Gateway,
		MacAddress:   s.MacAddress,
		Version:      s.Version,
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
