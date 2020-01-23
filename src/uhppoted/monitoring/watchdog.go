package monitoring

import (
	"fmt"
	"log"
	"math"
	"time"
	"uhppote/types"
)

const (
	// IDLE    = time.Duration(60 * time.Second)
	// IGNORE  = time.Duration(5 * time.Minute)
	// DELTA   = 60
	DELAY = 30
)

type Watchdog struct {
	healthcheck *HealthCheck
	log         *log.Logger
	state       struct {
		Started     time.Time
		HealthCheck struct {
			Alerted bool
		}

		// 	Devices struct {
		// 		Status sync.Map
		// 		Errors sync.Map
		// 	}
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
	w.log.Printf("INFO  %-20s", "watchdog")
	warnings := 0
	errors := 0
	healthCheckRunning := false
	// 	now := time.Now()

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
			if err := handler.Alert(w, msg); err != nil {
				w.log.Printf("WARN  %-20s %v", "monitoring", err)
			} else {
				w.state.HealthCheck.Alerted = true
			}
		}
	} else {
		if w.state.HealthCheck.Alerted {
			w.log.Printf("INFO  'health-check' subsystem is running")
			w.state.HealthCheck.Alerted = false
		}
	}

	// Verify configured devices

	if healthCheckRunning {
		// 		for id, _ := range u.Devices {
		// 			alerted := alerts{
		// 				missing:      false,
		// 				unexpected:   false,
		// 				touched:      false,
		// 				synchronized: false,
		// 			}

		// 			if v, found := st.devices.errors.Load(id); found {
		// 				alerted.missing = v.(alerts).missing
		// 				alerted.unexpected = v.(alerts).unexpected
		// 				alerted.touched = v.(alerts).touched
		// 				alerted.synchronized = v.(alerts).synchronized
		// 			}

		// 			if _, found := st.devices.status.Load(id); !found {
		// 				errors += 1
		// 				if !alerted.missing {
		// 					l.Printf("ERROR UTC0311-L0x %s device not found", types.SerialNumber(id))
		// 					alerted.missing = true
		// 				}
		// 			}

		// 			if v, found := st.devices.status.Load(id); found {
		// 				touched := v.(status).touched
		// 				t := time.Time(v.(status).status.SystemDateTime)
		// 				dt := time.Since(t).Round(seconds)
		// 				dtt := int64(math.Abs(time.Since(touched).Seconds()))

		// 				if alerted.missing {
		// 					l.Printf("ERROR UTC0311-L0x %s present", types.SerialNumber(id))
		// 					alerted.missing = false
		// 				}

		// 				if now.After(touched.Add(IDLE)) {
		// 					errors += 1
		// 					if !alerted.touched {
		// 						l.Printf("ERROR UTC0311-L0x %s no response for %s", types.SerialNumber(id), time.Since(touched).Round(seconds))
		// 						alerted.touched = true
		// 						alerted.synchronized = false
		// 					}
		// 				} else {
		// 					if alerted.touched {
		// 						l.Printf("INFO  UTC0311-L0x %s connected", types.SerialNumber(id))
		// 						alerted.touched = false
		// 					}
		// 				}

		// 				if dtt < DELTA/2 {
		// 					if int64(math.Abs(dt.Seconds())) > DELTA {
		// 						errors += 1
		// 						if !alerted.synchronized {
		// 							l.Printf("ERROR UTC0311-L0x %s system time not synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
		// 							alerted.synchronized = true
		// 						}
		// 					} else {
		// 						if alerted.synchronized {
		// 							l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", types.SerialNumber(id), types.DateTime(t), dt)
		// 							alerted.synchronized = false
		// 						}
		// 					}
		// 				}
		// 			}

		// 			st.devices.errors.Store(id, alerted)
		// 		}
	}

	// 	// Any unexpected devices?

	// 	st.devices.status.Range(func(key, value interface{}) bool {
	// 		alerted := alerts{
	// 			missing:      false,
	// 			unexpected:   false,
	// 			touched:      false,
	// 			synchronized: false,
	// 		}

	// 		if v, found := st.devices.errors.Load(key); found {
	// 			alerted.missing = v.(alerts).missing
	// 			alerted.unexpected = v.(alerts).unexpected
	// 			alerted.touched = v.(alerts).touched
	// 			alerted.synchronized = v.(alerts).synchronized
	// 		}

	// 		for id, _ := range u.Devices {
	// 			if id == key {
	// 				if alerted.unexpected {
	// 					l.Printf("ERROR UTC0311-L0x %s added to configuration", types.SerialNumber(key.(uint32)))
	// 					alerted.unexpected = false
	// 					st.devices.errors.Store(id, alerted)
	// 				}

	// 				return true
	// 			}
	// 		}

	// 		touched := value.(status).touched
	// 		t := time.Time(value.(status).status.SystemDateTime)
	// 		dt := time.Since(t).Round(seconds)
	// 		dtt := int64(math.Abs(time.Since(touched).Seconds()))

	// 		if now.After(touched.Add(IGNORE)) {
	// 			st.devices.status.Delete(key)
	// 			st.devices.errors.Delete(key)

	// 			if alerted.unexpected {
	// 				l.Printf("WARN  UTC0311-L0x %s disappeared", types.SerialNumber(key.(uint32)))
	// 			}
	// 		} else {
	// 			warnings += 1
	// 			if !alerted.unexpected {
	// 				l.Printf("WARN  UTC0311-L0x %s unexpected device", types.SerialNumber(key.(uint32)))
	// 				alerted.unexpected = true
	// 			}

	// 			if now.After(touched.Add(IDLE)) {
	// 				warnings += 1
	// 				if !alerted.touched {
	// 					l.Printf("WARN  UTC0311-L0x %s no response for %s", types.SerialNumber(key.(uint32)), time.Since(touched).Round(seconds))
	// 					alerted.touched = true
	// 					alerted.synchronized = false
	// 				}
	// 			} else {
	// 				if alerted.touched {
	// 					l.Printf("INFO  UTC0311-L0x %s connected", types.SerialNumber(key.(uint32)))
	// 					alerted.touched = false
	// 				}
	// 			}

	// 			if dtt < DELTA/2 {
	// 				if int64(math.Abs(dt.Seconds())) > DELTA {
	// 					warnings += 1
	// 					if !alerted.synchronized {
	// 						l.Printf("WARN  UTC0311-L0x %s system time not synchronized: %s (%s)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
	// 						alerted.synchronized = true
	// 					}
	// 				} else {
	// 					if alerted.synchronized {
	// 						l.Printf("INFO   UTC0311-L0x %s system time synchronized: %s (%s)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
	// 						alerted.synchronized = false
	// 					}
	// 				}
	// 			}

	// 			st.devices.errors.Store(key, alerted)
	// 		}

	// 		return true
	// 	})

	// 'k, done

	if errors > 0 {
		w.log.Printf("ERROR watchdog found %d errors", errors)
	}

	if warnings > 0 {
		w.log.Printf("WARN  watchdog found %d warnings", warnings)
	}

	if errors == 0 && warnings == 0 {
		w.log.Printf("INFO  watchdog: OK")
	}

	if err := handler.Alive(w, "watchdog"); err != nil {
		w.log.Printf("WARN  %-20s %v", "monitoring", err)
	}

	return nil
}
