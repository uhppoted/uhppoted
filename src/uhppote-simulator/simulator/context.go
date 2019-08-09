package simulator

import (
	"fmt"
	"math/rand"
	"net"
	"path"
	"uhppote-simulator/entities"
	"uhppote/types"
)

type DeviceList struct {
	txq     chan entities.Message
	devices []*Simulator
}

type Context struct {
	DeviceList DeviceList
	Directory  string
}

func NewDeviceList(l []*Simulator) DeviceList {
	txq := make(chan entities.Message, 8)
	for _, s := range l {
		s.TxQ = txq
	}

	return DeviceList{
		txq:     txq,
		devices: l,
	}
}

func (l *DeviceList) GetMessage() entities.Message {
	return <-l.txq
}

func (l *DeviceList) Apply(f func(*Simulator)) {
	for _, s := range l.devices {
		f(s)
	}
}

func (l *DeviceList) Find(deviceId uint32) *Simulator {
	for _, s := range l.devices {
		if s.SerialNumber == types.SerialNumber(deviceId) {
			return s
		}
	}

	return nil
}

func (l *DeviceList) Add(deviceId uint32, compressed bool, dir string) error {
	for _, s := range l.devices {
		if s.SerialNumber == types.SerialNumber(deviceId) {
			return nil
		}
	}

	filename := fmt.Sprintf("%d.json", deviceId)
	if compressed {
		filename = fmt.Sprintf("%d.json.gz", deviceId)
	}

	mac := make([]byte, 6)
	rand.Read(mac)

	device := Simulator{
		File:         path.Join(dir, filename),
		Compressed:   compressed,
		TxQ:          l.txq,
		SerialNumber: types.SerialNumber(deviceId),
		IpAddress:    net.IPv4(0, 0, 0, 0),
		SubnetMask:   net.IPv4(255, 255, 255, 0),
		Gateway:      net.IPv4(0, 0, 0, 0),
		MacAddress:   types.MacAddress(mac),
		Version:      0x0892,
	}

	err := (&device).Save()
	if err != nil {
		return err
	}

	l.devices = append(l.devices, &device)

	return nil
}

func (l *DeviceList) Delete(deviceId uint32) error {
	for ix, s := range l.devices {
		if s.SerialNumber == types.SerialNumber(deviceId) {
			if err := s.Delete(); err != nil {
				return err
			}

			l.devices = append(l.devices[:ix], l.devices[ix+1:]...)
		}
	}

	return nil
}
