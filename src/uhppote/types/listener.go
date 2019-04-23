package types

import (
	"fmt"
	"net"
)

type Listener struct {
	SerialNumber SerialNumber
	Address      net.UDPAddr
}

func (l *Listener) String() string {
	return fmt.Sprintf("%s %s", l.SerialNumber, l.Address.String())
}
