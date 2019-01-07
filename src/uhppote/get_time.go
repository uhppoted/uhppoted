package uhppote

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func GetTime(serialNumber uint32, debug bool) (*types.DateTime, error) {
	reply := make([]byte, 2048)
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x32
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

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
		fmt.Printf(" ... received %v\n", N)
		fmt.Printf(" ... received %v\n", remote)
		fmt.Printf(" ... received\n%s\n", hex.Dump(reply[:N]))
	}

	result, err := messages.NewGetTime(reply)

	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf(" ... %v\n", *result)
	}

	return &result.DateTime, nil
}
