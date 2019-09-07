package uhppote

import (
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"testing"
	"uhppote"
	"uhppote/messages"
)

func TestBroadcastAddressRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := messages.DeleteCardRequest{
		SerialNumber: 423187757,
		CardNumber:   6154412,
	}

	bind, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	broadcast, _ := net.ResolveUDPAddr("udp", "127.0.0.1:60000")

	u := uhppote.UHPPOTE{
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            true,
		BindAddress:      bind,
		BroadcastAddress: broadcast,
	}

	closed := make(chan int)
	c := listen("127.0.0.1:60000", closed, t)

	if c != nil {
		defer func() {
			c.Close()
		}()
	}

	reply, err := u.Send(423187757, request)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if reply == nil {
		t.Fatalf("Invalid reply: %v", reply)
	}

	if !reflect.DeepEqual(reply, expected) {
		t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", dump(expected, ""), dump(reply, ""))
	}

	c.Close()

	<-closed
}

func TestDeviceDirectedRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := messages.DeleteCardRequest{
		SerialNumber: 423187757,
		CardNumber:   6154412,
	}

	bind, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	broadcast, _ := net.ResolveUDPAddr("udp", "127.0.0.1:60000")
	d423187757, _ := net.ResolveUDPAddr("udp", "127.0.0.1:60001")
	d757781324, _ := net.ResolveUDPAddr("udp", "127.0.0.1:60002")

	u := uhppote.UHPPOTE{
		Debug:            true,
		BindAddress:      bind,
		BroadcastAddress: broadcast,
		Devices: map[uint32]*net.UDPAddr{
			423187757: d423187757,
			757781324: d757781324,
		},
	}

	closed := make(chan int)
	c := listen("127.0.0.1: 60001", closed, t)

	if c != nil {
		defer func() {
			c.Close()
		}()
	}

	reply, err := u.Send(423187757, request)

	if err != nil {
		t.Fatalf("%v", err)
	}

	if reply == nil {
		t.Fatalf("Invalid reply: %v", reply)
	}

	if !reflect.DeepEqual(reply, expected) {
		t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", dump(expected, ""), dump(reply, ""))
	}

	c.Close()

	<-closed
}

func listen(address string, closed chan int, t *testing.T) *net.UDPConn {
	response := []byte{
		0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		t.Fatalf("Error setting up test UDP device: %v", err)
	}

	c, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("%v", err)
	}

	go func() {
		m := make([]byte, 2048)

		for {
			_, remote, err := c.ReadFromUDP(m)
			if err != nil {
				t.Logf("%v", err)
				break
			}

			_, err = c.WriteTo(response, remote)
			if err != nil {
				t.Logf("%v", err)
				break
			}
		}

		closed <- 0
	}()

	return c
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
