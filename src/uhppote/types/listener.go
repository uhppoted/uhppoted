package types

import (
	"fmt"
	"net"
)

type Listener struct {
	SerialNumber SerialNumber
	Address      net.UDPAddr
}

type ListenerResult struct {
	SerialNumber SerialNumber
	Address      net.UDPAddr
	Succeeded    bool
}

func (l *Listener) String() string {
	return fmt.Sprintf("%s %s", l.SerialNumber, l.Address.String())
}

func (l *ListenerResult) String() string {
	return fmt.Sprintf("%s %s %s", l.SerialNumber, l.Address.String(), l.Succeeded)
}
