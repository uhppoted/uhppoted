package uhppote

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"time"
	"uhppote/messages"
	"uhppote/types"
)

func SetTime(serialNumber uint32, datetime time.Time, debug bool) (*types.DateTime, error) {
	reply := make([]byte, 2048)
	cmd := make([]byte, 64)

	cmd[0] = 0x17
	cmd[1] = 0x30
	cmd[2] = 0x00
	cmd[3] = 0x00

	binary.LittleEndian.PutUint32(cmd[4:8], serialNumber)

	cmd[8] = encode(datetime.Year() / 100)
	cmd[9] = encode(datetime.Year() % 100)
	cmd[10] = encode(int(datetime.Month()))
	cmd[11] = encode(datetime.Day())
	cmd[12] = encode(datetime.Hour())
	cmd[13] = encode(datetime.Minute())
	cmd[14] = encode(datetime.Second())

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

	result, err := messages.NewGetTime(reply)

	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf(" ... %v\n", *result)
	}

	return &result.DateTime, nil
}

func encode(b int) byte {
	msb := b / 10
	lsb := b % 10

	return byte(msb*16 + lsb)
}
