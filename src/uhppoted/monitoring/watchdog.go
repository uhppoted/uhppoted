package monitoring

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"log"
	"math"
	"time"
)

type Watchdog struct {
	healthcheck *HealthCheck
	log         *log.Logger
	state       struct {
		Started     time.Time
		HealthCheck struct {
			Alerted bool
		}
	}
}

func NewWatchdog(h *HealthCheck, l *log.Logger) Watchdog {
	return Watchdog{
		healthcheck: h,
		log:         l,
		state: struct {
			Started     time.Time
			HealthCheck struct {
				Alerted bool
			}
		}{
			Started: time.Now(),
			HealthCheck: struct {
				Alerted bool
			}{
				Alerted: false,
			},
		},
	}
}

func (w *Watchdog) ID() string {
	return "watchdog"
}

func (w *Watchdog) Exec(handler MonitoringHandler) error {
	w.log.Printf("DEBUG %-20s", "watchdog")

	warnings := uint(0)
	errors := uint(0)
	healthCheckRunning := false

	// Verify health-check
	dt := time.Since(w.state.Started).Round(time.Second)
	if w.healthcheck.state.Touched != nil {
		dt = time.Since(*w.healthcheck.state.Touched)
		if int64(math.Abs(dt.Seconds())) < DELAY {
			healthCheckRunning = true
		}
	}

	if int64(math.Abs(dt.Seconds())) > DELAY {
		errors += 1
		if !w.state.HealthCheck.Alerted {
			msg := fmt.Sprintf("'health-check' subsystem has not run since %s (%s)", types.DateTime(w.state.Started), dt)

			w.log.Printf("ERROR %s", msg)
			if err := handler.Alert(w, msg); err == nil {
				w.state.HealthCheck.Alerted = true
			}
		}
	} else {
		if w.state.HealthCheck.Alerted {
			w.log.Printf("INFO  'health-check' subsystem is running")
			w.state.HealthCheck.Alerted = false
		}
	}

	// Report on known devices
	if healthCheckRunning {
		warnings += w.healthcheck.state.Warnings
		errors += w.healthcheck.state.Errors
	}

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

	w.log.Printf("%-5s %-12s %s", level, "watchdog", msg)
	handler.Alive(w, msg)

	return nil
}
