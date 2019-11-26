package simulator

import (
	"uhppote-simulator/entities"
	"uhppote-simulator/simulator/UT0311L04"
)

type DeviceList struct {
	txq     chan entities.Message
	devices []Simulator
}

type Context struct {
	BindAddress string
	DeviceList  DeviceList
	Directory   string
}

func NewDeviceList(l []Simulator) DeviceList {
	txq := make(chan entities.Message, 8)
	for _, s := range l {
		s.SetTxQ(txq)
	}

	return DeviceList{
		txq:     txq,
		devices: l,
	}
}

func (l *DeviceList) GetMessage() entities.Message {
	return <-l.txq
}

func (l *DeviceList) Apply(f func(Simulator)) {
	for _, s := range l.devices {
		f(s)
	}
}

func (l *DeviceList) Find(deviceId uint32) Simulator {
	for _, s := range l.devices {
		if s.DeviceID() == deviceId {
			return s
		}
	}

	return nil
}

func (l *DeviceList) Add(deviceId uint32, compressed bool, dir string) (bool, error) {
	for _, s := range l.devices {
		if s.DeviceID() == deviceId {
			return false, nil
		}
	}

	device := UT0311L04.NewUT0311L04(deviceId, dir, compressed)
	device.SetTxQ(l.txq)
	err := device.Save()
	if err != nil {
		return false, err
	}

	l.devices = append(l.devices, device)

	return true, nil
}

func (l *DeviceList) Delete(deviceId uint32) error {
	for ix, s := range l.devices {
		if s.DeviceID() == deviceId {
			if err := s.Delete(); err != nil {
				return err
			}

			l.devices = append(l.devices[:ix], l.devices[ix+1:]...)
		}
	}

	return nil
}
