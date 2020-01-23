package monitoring

import ()

type Monitor interface {
	ID() string
}

type MonitoringHandler interface {
	Alive(Monitor, string) error
	Alert(Monitor, string) error
}
