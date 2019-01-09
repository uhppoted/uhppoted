package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"time"
)

func Execute(cmd []byte, debug bool) ([]byte, error) {
	reply := make([]byte, 2048)

	if debug {
		fmt.Printf(" ... sent:\n%s\n", hex.Dump(cmd))
	}

	local, err := net.ResolveUDPAddr("udp", "192.168.1.100:50000")

	if err != nil {
		return nil, makeErr("Failed to resolve UDP local address", err)
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return nil, makeErr("Failed to resolve UDP broadcast address", err)
	}

	connection, err := net.ListenUDP("udp", local)

	if err != nil {
		return nil, makeErr("Failed to open UDP socket", err)
	}

	defer close(connection)

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return nil, makeErr("Failed to write to UDP socket", err)
	}

	if debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

	err = connection.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	if err != nil {
		return nil, makeErr("Failed to set UDP timeout", err)
	}

	N, remote, err := connection.ReadFromUDP(reply)

	if err != nil {
		return nil, makeErr("Failed to read from UDP socket", err)
	}

	if debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... received %v bytes from %v\n", N, remote)
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ... $1"))
	}

	return reply[:N], nil
}
