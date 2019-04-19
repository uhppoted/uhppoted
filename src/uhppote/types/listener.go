package types

import (
	"fmt"
	"net"
)

type Listener struct {
	SerialNumber SerialNumber
	Address      net.IP
	Port         uint16
}

func (l *Listener) String() string {
	return fmt.Sprintf("%s %v:%d",
		l.SerialNumber,
		l.Address,
		l.Port)
}
