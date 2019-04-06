package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"regexp"
	"time"
	"uhppote/encoding"
	"uhppote/messages"
)

type UHPPOTE struct {
	BindAddress net.UDPAddr
	Debug       bool
}

func (u *UHPPOTE) Exec(m messages.Message) ([]byte, error) {
	request, err := uhppote.Marshal(m)

	if err != nil {
		return nil, makeErr(fmt.Sprintf("Error encoding request %v", m), err)
	}

	return u.Execute(request)
}

func (u *UHPPOTE) Execute(cmd []byte) ([]byte, error) {
	reply := make([]byte, 2048)

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... command %v bytes\n", len(cmd))
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(cmd), " ...         $1"))
	}

	broadcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:60000")

	if err != nil {
		return nil, makeErr("Failed to resolve UDP broadcast address", err)
	}

	connection, err := net.ListenUDP("udp", &u.BindAddress)

	if err != nil {
		return nil, makeErr("Failed to open UDP socket", err)
	}

	defer close(connection)

	N, err := connection.WriteTo(cmd, broadcast)

	if err != nil {
		return nil, makeErr("Failed to write to UDP socket", err)
	}

	if u.Debug {
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

	if u.Debug {
		regex := regexp.MustCompile("(?m)^(.*)")

		fmt.Printf(" ... received %v bytes from %v\n", N, remote)
		fmt.Printf("%s\n", regex.ReplaceAllString(hex.Dump(reply[:N]), " ...          $1"))
	}

	return reply[:N], nil
}

func close(connection net.Conn) {
	fmt.Println(" ... closing connection")

	connection.Close()
}
