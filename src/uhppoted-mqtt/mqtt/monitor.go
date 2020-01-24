package mqtt

import (
	"sync"
	"time"
	"uhppoted/monitoring"
)

type SystemMonitor struct {
	mqttd *MQTTD
}

var alive = sync.Map{}

func NewSystemMonitor(mqttd *MQTTD) *SystemMonitor {
	return &SystemMonitor{
		mqttd: mqttd,
	}
}

func (m *SystemMonitor) Alive(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alive string `json:"alive"`
	}{
		Alive: msg,
	}

	now := time.Now()
	last, ok := alive.Load(monitor.ID())
	interval := 60 * time.Second

	if ok && time.Since(last.(time.Time)).Round(time.Second) < interval {
		return nil
	}

	if err := m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, event, msgSystem); err != nil {
		return err
	}

	alive.Store(monitor.ID(), now)

	return nil
}

func (m *SystemMonitor) Alert(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alert string `json:"alert"`
	}{
		Alert: msg,
	}

	return m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, event, msgSystem)
}
