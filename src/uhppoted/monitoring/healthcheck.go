package monitoring

import (
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"github.com/uhppoted/uhppoted/src/uhppote/types"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

type HealthCheck struct {
	uhppote    *uhppote.UHPPOTE
	idleTime   time.Duration
	ignoreTime time.Duration
	log        *log.Logger
	state      struct {
		Started time.Time
		Touched *time.Time
		Devices struct {
			Status   sync.Map
			Listener sync.Map
			Errors   sync.Map
		}
		Warnings uint
		Errors   uint
	}
}

type status struct {
	Touched time.Time
	Status  types.Status
}

type listener struct {
	Touched time.Time
	Address net.UDPAddr
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
	nolistener   bool
	listener     bool
}

func NewHealthCheck(u *uhppote.UHPPOTE, idleTime, ignoreTime time.Duration, l *log.Logger) HealthCheck {
	return HealthCheck{
		uhppote:    u,
		idleTime:   idleTime,
		ignoreTime: ignoreTime,
		log:        l,
		state: struct {
			Started time.Time
			Touched *time.Time
			Devices struct {
				Status   sync.Map
				Listener sync.Map
				Errors   sync.Map
			}
			Warnings uint
			Errors   uint
		}{
			Started: time.Now(),
			Touched: nil,
			Devices: struct {
				Status   sync.Map
				Listener sync.Map
				Errors   sync.Map
			}{
				Status:   sync.Map{},
				Listener: sync.Map{},
				Errors:   sync.Map{},
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
		msg = fmt.Sprintf("%s", Errors(errors))
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
				Status:  *s,
				Touched: now,
			})
		}

		l, err := h.uhppote.GetListener(id)
		if err == nil && l != nil {
			h.state.Devices.Listener.Store(id, listener{
				Address: l.Address,
				Touched: now,
			})
		}
	}

	h.state.Touched = &now
}

// Check known/identified devices
func (h *HealthCheck) known(now time.Time, handler MonitoringHandler) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	for id, _ := range h.uhppote.Devices {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
			nolistener:   false,
			listener:     false,
		}

		if v, found := h.state.Devices.Errors.Load(id); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
			alerted.nolistener = v.(alerts).nolistener
			alerted.listener = v.(alerts).listener
		}

		if _, found := h.state.Devices.Status.Load(id); !found {
			errors += 1
			if !alerted.missing {
				if alert(h, handler, id, "device not found") {
					alerted.missing = true
				}
			}
		} else {
			if alerted.missing {
				if info(h, handler, id, "device present") {
					alerted.missing = false
				}
			}
		}

		e, w := h.checkStatus(id, now, &alerted, handler, true)
		errors += e
		warnings += w

		e, w = h.checkListener(id, now, &alerted, handler, true)
		errors += e
		warnings += w

		h.state.Devices.Errors.Store(id, alerted)
	}

	return errors, warnings
}

// Identify and check any unexpected devices
func (h *HealthCheck) unexpected(now time.Time, handler MonitoringHandler) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

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
			alerted.nolistener = v.(alerts).nolistener
			alerted.listener = v.(alerts).listener
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
		if now.After(touched.Add(h.ignoreTime)) {
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

			e, w := h.checkStatus(key.(uint32), now, &alerted, handler, false)
			errors += e
			warnings += w

			e, w = h.checkListener(key.(uint32), now, &alerted, handler, false)
			errors += e
			warnings += w

			h.state.Devices.Errors.Store(key, alerted)
		}

		return true
	}

	h.state.Devices.Status.Range(f)

	return errors, warnings
}

func (h *HealthCheck) checkStatus(id uint32, now time.Time, alerted *alerts, handler MonitoringHandler, known bool) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	if v, found := h.state.Devices.Status.Load(id); found {
		touched := v.(status).Touched
		t := time.Time(v.(status).Status.SystemDateTime)
		dt := time.Since(t).Round(time.Second)
		dtt := int64(math.Abs(time.Since(touched).Seconds()))

		if now.After(touched.Add(h.idleTime)) {
			if known {
				errors += 1
			} else {
				warnings += 1
			}

			if !alerted.touched {
				msg := fmt.Sprintf("no response for %s", time.Since(touched).Round(time.Second))
				if alert(h, handler, id, msg) {
					alerted.touched = true
					alerted.synchronized = false
				}
			}
		} else {
			if alerted.touched {
				if info(h, handler, id, "connected") {
					alerted.touched = false
				}
			}
		}

		if dtt < DELTA/2 {
			if int64(math.Abs(dt.Seconds())) > DELTA {
				if known {
					errors += 1
				} else {
					warnings += 1
				}

				if !alerted.synchronized {
					msg := fmt.Sprintf("system time not synchronized: %s (%s)", types.DateTime(t), dt)
					if alert(h, handler, id, msg) {
						alerted.synchronized = true
					}
				}
			} else {
				if alerted.synchronized {
					msg := fmt.Sprintf("system time synchronized: %s (%s)", types.DateTime(t), dt)
					if info(h, handler, id, msg) {
						alerted.synchronized = false
					}
				}
			}
		}
	}

	return errors, warnings
}

func (h *HealthCheck) checkListener(id uint32, now time.Time, alerted *alerts, handler MonitoringHandler, known bool) (uint, uint) {
	errors := uint(0)
	warnings := uint(0)

	expected := h.uhppote.ListenAddress
	if expected == nil {
		return errors, warnings
	}

	if v, found := h.state.Devices.Listener.Load(id); found {
		address := v.(listener).Address
		touched := v.(listener).Touched

		if now.After(touched.Add(h.idleTime)) {
			if known {
				errors += 1
			} else {
				warnings += 1
			}

			if !alerted.nolistener {
				msg := fmt.Sprintf("no reply to 'get-listener' for %s", time.Since(touched).Round(time.Second))
				if warn(h, handler, id, msg) {
					alerted.nolistener = true
				}
			}
		} else {
			if alerted.nolistener {
				if info(h, handler, id, "listener identified") {
					alerted.nolistener = false
				}
			}
		}

		if !expected.IP.Equal(address.IP) || (expected.Port != address.Port) {
			if known {
				errors += 1
			} else {
				warnings += 1
			}

			if !alerted.listener {
				msg := fmt.Sprintf("incorrect listener address/port: %s", &address)
				if warn(h, handler, id, msg) {
					alerted.listener = true
				}
			}
		} else {
			if alerted.listener {
				if info(h, handler, id, "listener address/port correct") {
					alerted.listener = false
				}
			}
		}
	}

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
	known := false

	for id, _ := range h.uhppote.Devices {
		if deviceID == id {
			known = true
		}
	}

	if known {
		h.log.Printf("%-5s %s", "ERROR", msg)
	} else {
		h.log.Printf("%-5s %s", "WARN", msg)
	}

	if err := handler.Alert(h, msg); err != nil {
		return false
	}

	return true
}
