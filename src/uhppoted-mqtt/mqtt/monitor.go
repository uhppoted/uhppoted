package mqtt

import (
	"github.com/uhppoted/uhppoted-api/monitoring"
	"log"
	"sync"
	"time"
)

type SystemMonitor struct {
	mqttd *MQTTD
	log   *log.Logger
}

var alive = sync.Map{}

func NewSystemMonitor(mqttd *MQTTD, log *log.Logger) *SystemMonitor {
	return &SystemMonitor{
		mqttd: mqttd,
		log:   log,
	}
}

func (m *SystemMonitor) Alive(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alive struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		} `json:"alive"`
	}{
		Alive: struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		}{
			SubSystem: monitor.ID(),
			Message:   msg,
		},
	}

	now := time.Now()
	last, ok := alive.Load(monitor.ID())
	interval := 60 * time.Second

	if ok && time.Since(last.(time.Time)).Round(time.Second) < interval {
		return nil
	}

	if err := m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, event, msgSystem, false); err != nil {
		m.log.Printf("WARN  %-20s %v", "monitoring", err)
		return err
	}

	alive.Store(monitor.ID(), now)

	return nil
}

func (m *SystemMonitor) Alert(monitor monitoring.Monitor, msg string) error {
	event := struct {
		Alert struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		} `json:"alert"`
	}{
		Alert: struct {
			SubSystem string `json:"subsystem"`
			Message   string `json:"message"`
		}{
			SubSystem: monitor.ID(),
			Message:   msg,
		},
	}

	if err := m.mqttd.send(&m.mqttd.Encryption.SystemKeyID, m.mqttd.Topics.System, event, msgSystem, true); err != nil {
		m.log.Printf("WARN  %-20s %v", "monitoring", err)
		return err
	}

	return nil
}
