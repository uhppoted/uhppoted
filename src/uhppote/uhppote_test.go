package uhppote

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
	"uhppote"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
	"uhppote/types"
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

	u := uhppote.UHPPOTE{
		Devices:          make(map[uint32]*net.UDPAddr),
		Debug:            true,
		BindAddress:      *resolve("127.0.0.1:12345", t),
		BroadcastAddress: *resolve("127.0.0.1:60000", t),
	}

	closed := make(chan int)
	c := listen(423187757, "127.0.0.1:60000", 0*time.Millisecond, closed, t)

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
		t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected, ""), uhppote.Dump(reply, ""))
	}

	c.Close()

	<-closed
}

func TestSequentialRequests(t *testing.T) {
	expected := [][]byte{
		{0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		{0x17, 0x52, 0x00, 0x00, 0x4c, 0xd3, 0x2a, 0x2d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	}

	request := messages.DeleteCardRequest{
		SerialNumber: 1000002,
		CardNumber:   6154412,
	}

	u := uhppote.UHPPOTE{
		Debug:            true,
		BindAddress:      *resolve("127.0.0.1:12345", t),
		BroadcastAddress: *resolve("127.0.0.1:60000", t),
		Devices: map[uint32]*net.UDPAddr{
			423187757: resolve("127.0.0.1:60001", t),
			757781324: resolve("127.0.0.1:60002", t),
		},
	}

	closed := make(chan int)
	listening := []*net.UDPConn{
		listen(423187757, "127.0.0.1: 60001", 0*time.Millisecond, closed, t),
		listen(757781324, "127.0.0.1: 60002", 0*time.Millisecond, closed, t),
	}

	defer func() {
		for _, c := range listening {
			if c != nil {
				c.Close()
			}
		}
	}()

	if reply, err := u.Send(423187757, request); err != nil {
		t.Fatalf("%v", err)
	} else if reply == nil {
		t.Fatalf("Invalid reply: %v", reply)
	} else if !reflect.DeepEqual(reply, expected[0]) {
		t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[0], ""), uhppote.Dump(reply, ""))
	}

	if reply, err := u.Send(757781324, request); err != nil {
		t.Fatalf("%v", err)
	} else if reply == nil {
		t.Fatalf("Invalid reply: %v", reply)
	} else if !reflect.DeepEqual(reply, expected[1]) {
		t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[1], ""), uhppote.Dump(reply, ""))
	}

	for _, c := range listening {
		if c != nil {
			c.Close()
			<-closed
		}
	}
}

func TestConcurrentRequestsWithUnboundPort(t *testing.T) {
	expected := [][]byte{
		{0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		{0x17, 0x52, 0x00, 0x00, 0x4c, 0xd3, 0x2a, 0x2d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	}

	request := messages.DeleteCardRequest{
		SerialNumber: 1000002,
		CardNumber:   6154412,
	}

	u := uhppote.UHPPOTE{
		Debug:            true,
		BindAddress:      *resolve("127.0.0.1:0", t),
		BroadcastAddress: *resolve("127.0.0.1:60000", t),
		Devices: map[uint32]*net.UDPAddr{
			423187757: resolve("127.0.0.1:60001", t),
			757781324: resolve("127.0.0.1:60002", t),
		},
	}

	closed := make(chan int)
	listening := []*net.UDPConn{
		listen(423187757, "127.0.0.1: 60001", 1500*time.Millisecond, closed, t),
		listen(757781324, "127.0.0.1: 60002", 500*time.Millisecond, closed, t),
	}

	defer func() {
		for _, c := range listening {
			if c != nil {
				c.Close()
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		if reply, err := u.Send(423187757, request); err != nil {
			t.Fatalf("%v", err)
		} else if reply == nil {
			t.Fatalf("Invalid reply: %v", reply)
		} else if !reflect.DeepEqual(reply, expected[0]) {
			t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[0], ""), uhppote.Dump(reply, ""))
		}
	}()

	go func() {
		defer wg.Done()

		time.Sleep(500 * time.Millisecond)

		if reply, err := u.Send(757781324, request); err != nil {
			t.Fatalf("%v", err)
		} else if reply == nil {
			t.Fatalf("Invalid reply: %v", reply)
		} else if !reflect.DeepEqual(reply, expected[1]) {
			t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[1], ""), uhppote.Dump(reply, ""))
		}
	}()

	wg.Wait()

	for _, c := range listening {
		if c != nil {
			c.Close()
			<-closed
		}
	}
}

func TestConcurrentRequestsWithBoundPort(t *testing.T) {
	t.Skip("SKIP - uhppote concurrency implementation is a work in progress")

	expected := [][]byte{
		{0x17, 0x52, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
		{0x17, 0x52, 0x00, 0x00, 0x4c, 0xd3, 0x2a, 0x2d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		},
	}

	request := messages.DeleteCardRequest{
		SerialNumber: 1000002,
		CardNumber:   6154412,
	}

	u := uhppote.UHPPOTE{
		Debug:            true,
		BindAddress:      *resolve("127.0.0.1:12345", t),
		BroadcastAddress: *resolve("127.0.0.1:60000", t),
		Devices: map[uint32]*net.UDPAddr{
			423187757: resolve("127.0.0.1:60001", t),
			757781324: resolve("127.0.0.1:60002", t),
		},
	}

	closed := make(chan int)
	listening := []*net.UDPConn{
		listen(423187757, "127.0.0.1: 60001", 1000*time.Millisecond, closed, t),
		listen(757781324, "127.0.0.1: 60002", 500*time.Millisecond, closed, t),
	}

	defer func() {
		for _, c := range listening {
			if c != nil {
				c.Close()
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		if reply, err := u.Send(423187757, request); err != nil {
			t.Fatalf("%v", err)
		} else if reply == nil {
			t.Fatalf("Invalid reply: %v", reply)
		} else if !reflect.DeepEqual(reply, expected[0]) {
			t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[0], ""), uhppote.Dump(reply, ""))
		}
	}()

	go func() {
		defer wg.Done()

		time.Sleep(500 * time.Millisecond)

		if reply, err := u.Send(757781324, request); err != nil {
			t.Fatalf("%v", err)
		} else if reply == nil {
			t.Fatalf("Invalid reply: %v", reply)
		} else if !reflect.DeepEqual(reply, expected[1]) {
			t.Fatalf("Incorrect reply:\nExpected:\n%s\nReturned:\n%s", uhppote.Dump(expected[1], ""), uhppote.Dump(reply, ""))
		}
	}()

	wg.Wait()

	for _, c := range listening {
		if c != nil {
			c.Close()
			<-closed
		}
	}
}

func listen(deviceId uint32, address string, delay time.Duration, closed chan int, t *testing.T) *net.UDPConn {
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

			response := messages.DeleteCardResponse{
				SerialNumber: types.SerialNumber(deviceId),
				Succeeded:    true,
			}

			m, err := codec.Marshal(response)
			if err != nil {
				t.Logf("%v", err)
				break
			}

			time.Sleep(delay)

			_, err = c.WriteTo(m, remote)
			if err != nil {
				t.Logf("%v", err)
				break
			}
		}

		closed <- 1
	}()

	return c
}

func resolve(address string, t *testing.T) *net.UDPAddr {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		t.Fatalf("Error resolving UDP address '%s': %v", address, err)
	}

	return addr
}
