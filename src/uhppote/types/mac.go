package types

import (
	"encoding/json"
	"net"
)

type MacAddress net.HardwareAddr

func (m MacAddress) String() string {
	return net.HardwareAddr(m).String()
}

func (m MacAddress) MarshalUT0311L0x() ([]byte, error) {
	bytes := make([]byte, 6)

	copy(bytes, m)

	return bytes, nil
}

func (m *MacAddress) UnmarshalUT0311L0x(bytes []byte) error {
	mac := make([]byte, 6)

	copy(mac, bytes[0:6])

	*m = MacAddress(mac)

	return nil
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
