package monitoring

import (
	"log"
	"sync"
	"time"
	"uhppote"
	"uhppote/types"
)

type HealthCheck struct {
	uhppote *uhppote.UHPPOTE
	log     *log.Logger
	state   struct {
		Started time.Time
		Touched *time.Time
		Devices struct {
			Status sync.Map
			Errors sync.Map
		}
	}
}

type status struct {
	Touched time.Time
	Status  types.Status
}

func NewHealthCheck(u *uhppote.UHPPOTE, l *log.Logger) HealthCheck {
	return HealthCheck{
		uhppote: u,
		log:     l,
		state: struct {
			Started time.Time
			Touched *time.Time
			Devices struct {
				Status sync.Map
				Errors sync.Map
			}
		}{
			Started: time.Now(),
			Touched: nil,
			Devices: struct {
				Status sync.Map
				Errors sync.Map
			}{
				Status: sync.Map{},
				Errors: sync.Map{},
			},
		},
	}
}

func (h *HealthCheck) ID() string {
	return "health-check"
}

func (h *HealthCheck) Exec(handler MonitoringHandler) {
	h.log.Printf("INFO  %-20s", "health-check")

	h.update()

	if err := handler.Alive(h, "health-check"); err != nil {
		h.log.Printf("WARN  %-20s %v", "monitoring", err)
	}
}

func (h *HealthCheck) update() {
	now := time.Now()
	devices := make(map[uint32]bool)

	found, err := h.uhppote.FindDevices()
	if err != nil {
		h.log.Printf("WARN  'keep-alive' error: %v", err)
	}

	if found != nil {
		for _, id := range found {
			devices[uint32(id.SerialNumber)] = true
		}
	}

	for id, _ := range h.uhppote.Devices {
		devices[id] = true
	}

	for id, _ := range devices {
		s, err := h.uhppote.GetStatus(id)
		if err == nil {
			h.state.Devices.Status.Store(id, status{
				Touched: now,
				Status:  *s,
			})
		}
	}

	h.state.Touched = &now
}
