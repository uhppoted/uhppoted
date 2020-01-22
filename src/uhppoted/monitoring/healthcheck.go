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
	state   state
}

type state struct {
	Started time.Time

	HealthCheck struct {
		Touched *time.Time
		Alerted bool
	}

	Devices struct {
		Status sync.Map
		Errors sync.Map
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
		state: state{
			Started: time.Now(),

			HealthCheck: struct {
				Touched *time.Time
				Alerted bool
			}{
				Touched: nil,
				Alerted: false,
			},
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

func (h *HealthCheck) Exec() {
	h.log.Printf("INFO  %-20s", "health-check")

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

	h.state.HealthCheck.Touched = &now
}
