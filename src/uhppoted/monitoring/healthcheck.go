package monitoring

import (
	"fmt"
	"log"
	"math"
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
		Warnings uint
		Errors   uint
	}
}

type status struct {
	Touched time.Time
	Status  types.Status
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
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
			Warnings uint
			Errors   uint
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
			Warnings: 0,
			Errors:   0,
		},
	}
}

func (h *HealthCheck) ID() string {
	return "health-check"
}

func (h *HealthCheck) Exec(handler MonitoringHandler) {
	h.log.Printf("DEBUG %-20s", "health-check")

	now := time.Now()
	errors := uint(0)
	warnings := uint(0)

	h.update(now)

	e, w := h.known(now, handler)
	errors += e
	warnings += w

	e, w = h.unexpected(now, handler)
	errors += e
	warnings += w

	h.state.Warnings = warnings
	h.state.Errors = errors

	// 'k, done

	level := "INFO"
	msg := "OK"

	if errors > 0 && warnings > 0 {
		level = "WARN"
		msg = fmt.Sprintf("%s, %s", Errors(errors), Warnings(warnings))
	} else if errors > 0 {
		level = "WARN"
		msg = fmt.Sprintf("%s", Errors(warnings))
	} else if warnings > 0 {
		level = "WARN"
		msg = fmt.Sprintf("%s", Warnings(warnings))
	}

	h.log.Printf("%-5s %-12s %s", level, "health-check", msg)
	handler.Alive(h, msg)
}

func (h *HealthCheck) update(now time.Time) {
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

// Check known/identified devices
func (h *HealthCheck) known(now time.Time, handler MonitoringHandler) (uint, uint) {
	warnings := uint(0)
	errors := uint(0)

	for id, _ := range h.uhppote.Devices {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
		}

		if v, found := h.state.Devices.Errors.Load(id); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		if _, found := h.state.Devices.Status.Load(id); !found {
			errors += 1
			if !alerted.missing {
				if alert(h, handler, id, "device not found") {
					alerted.missing = true
				}
			}
		}

		if v, found := h.state.Devices.Status.Load(id); found {
			touched := v.(status).Touched
			t := time.Time(v.(status).Status.SystemDateTime)
			dt := time.Since(t).Round(time.Second)
			dtt := int64(math.Abs(time.Since(touched).Seconds()))

			if alerted.missing {
				if info(h, handler, id, "device present") {
					alerted.missing = false
				}
			}

			if now.After(touched.Add(IDLE)) {
				errors += 1
				if !alerted.touched {
					msg := fmt.Sprintf("no response for %s", time.Since(touched).Round(time.Second))
					if alert(h, handler, id, msg) {
						alerted.touched = true
						alerted.synchronized = false
					}
				}
			} else {
				if alerted.touched {
					if alert(h, handler, id, "connected") {
						alerted.touched = false
					}
				}
			}

			if dtt < DELTA/2 {
				if int64(math.Abs(dt.Seconds())) > DELTA {
					errors += 1
					if !alerted.synchronized {
						msg := fmt.Sprintf("system time not synchronized: %s (%s)", types.DateTime(t), dt)
						if alert(h, handler, id, msg) {
							alerted.synchronized = true
						}
					}
				} else {
					if alerted.synchronized {
						msg := fmt.Sprintf("system time synchronized: %s (%s)", types.DateTime(t), dt)
						if alert(h, handler, id, msg) {
							alerted.synchronized = false
						}
					}
				}
			}
		}

		h.state.Devices.Errors.Store(id, alerted)
	}

	return errors, warnings
}

// Identify and check any unexpected devices
func (h *HealthCheck) unexpected(now time.Time, handler MonitoringHandler) (uint, uint) {
	warnings := uint(0)
	errors := uint(0)

	f := func(key, value interface{}) bool {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
		}

		if v, found := h.state.Devices.Errors.Load(key); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		for id, _ := range h.uhppote.Devices {
			if id == key {
				if alerted.unexpected {
					if alert(h, handler, key.(uint32), "added to configuration") {
						alerted.unexpected = false
						h.state.Devices.Errors.Store(id, alerted)
					}
				}

				return true
			}
		}

		touched := value.(status).Touched
		t := time.Time(value.(status).Status.SystemDateTime)
		dt := time.Since(t).Round(time.Second)
		dtt := int64(math.Abs(time.Since(touched).Seconds()))

		if now.After(touched.Add(IGNORE)) {
			h.state.Devices.Status.Delete(key)
			h.state.Devices.Errors.Delete(key)

			if alerted.unexpected {
				warn(h, handler, key.(uint32), "disappeared")
			}
		} else {
			warnings += 1
			if !alerted.unexpected {
				if warn(h, handler, key.(uint32), "unexpected device") {
					alerted.unexpected = true
				}
			}

			if now.After(touched.Add(IDLE)) {
				warnings += 1
				if !alerted.touched {
					msg := fmt.Sprintf("no response for %s", time.Since(touched).Round(time.Second))
					if warn(h, handler, key.(uint32), msg) {
						alerted.touched = true
						alerted.synchronized = false
					}
				}
			} else {
				if alerted.touched {
					if info(h, handler, key.(uint32), "connected") {
						alerted.touched = false
					}
				}
			}

			if dtt < DELTA/2 {
				if int64(math.Abs(dt.Seconds())) > DELTA {
					warnings += 1
					if !alerted.synchronized {
						msg := fmt.Sprintf("system time not synchronized: %s (%s)", types.DateTime(t), dt)
						if warn(h, handler, key.(uint32), msg) {
							alerted.synchronized = true
						}
					}
				} else {
					if alerted.synchronized {
						msg := fmt.Sprintf("system time synchronized: %s (%s)", types.DateTime(t), dt)
						if info(h, handler, key.(uint32), msg) {
							alerted.synchronized = false
						}
					}
				}
			}

			h.state.Devices.Errors.Store(key, alerted)
		}

		return true
	}

	h.state.Devices.Status.Range(f)

	return errors, warnings
}

func info(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)

	h.log.Printf("%-5s %s", "INFO", msg)
	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}

func warn(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)

	h.log.Printf("%-5s %s", "WARN", msg)
	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}

func alert(h *HealthCheck, handler MonitoringHandler, deviceID uint32, message string) bool {
	msg := fmt.Sprintf("UTC0311-L0x %s %s", types.SerialNumber(deviceID), message)

	h.log.Printf("%-5s %s", "ERROR", msg)
	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}
