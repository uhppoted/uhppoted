package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func Search(debug bool) ([]types.Device, error) {
	devices := []types.Device{}
	reply := make([]byte, 2048)
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x94
	cmd[2] = 0x00
	cmd[3] = 0x00

	if debug {
		fmt.Printf(" ... sent:\n%s\n", hex.Dump(cmd))
	}

	local, err := net.ResolveUDPAddr("udp", "192.168.1.100:50000")

	if err != nil {
		return devices, makeErr("Failed to resolve UDP local address", err)
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return devices, makeErr("Failed to resolve UDP broadcast address", err)
	}

	connection, err := net.ListenUDP("udp", local)

	if err != nil {
		return devices, makeErr("Failed to open UDP socket", err)
	}

	defer close(connection)

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return devices, makeErr("Failed to write to UDP socket", err)
	}

	if debug {
		fmt.Printf(" ... sent %v bytes\n", N)
	}

	err = connection.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	if err != nil {
		return devices, makeErr("Failed to set UDP timeout", err)
	}

	N, remote, err := connection.ReadFromUDP(reply)

	if err != nil {
		return devices, makeErr("Failed to read from UDP socket", err)
	}

	if debug {
		fmt.Printf(" ... received %v\n", N)
		fmt.Printf(" ... received %v\n", remote)
		fmt.Printf(" ... received\n%s\n", hex.Dump(reply[:N]))
	}

	result, err := messages.NewSearch(reply)

	if err != nil {
		return devices, err
	}

	if debug {
		fmt.Printf(" ... %v\n", *result)
	}

	devices = append(devices, result.Device)

	return devices, nil
}
