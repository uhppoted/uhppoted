package monitoring

import (
	"fmt"
	"time"
)

type Monitor interface {
	ID() string
}

type MonitoringHandler interface {
	Alive(Monitor, string) error
	Alert(Monitor, string) error
}

type Errors uint

type Warnings uint

func (e Errors) String() string {
	if e == 1 {
		return fmt.Sprintf("%d error", uint(e))
	}
	return fmt.Sprintf("%d errors", uint(e))
}

func (w Warnings) String() string {
	if w == 1 {
		return fmt.Sprintf("%d warning", uint(w))
	}

	return fmt.Sprintf("%d warnings", uint(w))
}

const (
	IDLE   = time.Duration(60 * time.Second)
	IGNORE = time.Duration(5 * time.Minute)
	DELTA  = 60
	DELAY  = 30
)
