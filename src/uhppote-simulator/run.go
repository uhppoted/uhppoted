package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"uhppote"
	"uhppote-simulator/simulator"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
)

type handler struct {
	factory    func() messages.Request
	dispatcher func(*simulator.Simulator, messages.Request) (messages.Response, error)
}

var handlers = map[byte]*handler{
	0x20: &handler{
		nil,
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetStatus(rq.(*messages.GetStatusRequest))
		},
	},

	0x30: &handler{
		func() messages.Request { return new(uhppote.SetTimeRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.SetTime(rq.(*uhppote.SetTimeRequest))
		},
	},

	0x32: &handler{
		func() messages.Request { return new(uhppote.GetTimeRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetTime(rq.(*uhppote.GetTimeRequest))
		},
	},

	0x40: &handler{
		func() messages.Request { return new(uhppote.OpenDoorRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.OpenDoor(rq.(*uhppote.OpenDoorRequest))
		},
	},

	0x50: &handler{
		func() messages.Request { return new(uhppote.PutCardRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.PutCard(rq.(*uhppote.PutCardRequest))
		},
	},

	0x52: &handler{
		func() messages.Request { return new(uhppote.DeleteCardRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.DeleteCard(rq.(*uhppote.DeleteCardRequest))
		},
	},

	0x54: &handler{
		func() messages.Request { return new(uhppote.DeleteCardsRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.DeleteCards(rq.(*uhppote.DeleteCardsRequest))
		},
	},

	0x58: &handler{
		func() messages.Request { return new(uhppote.GetCardsRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetCards(rq.(*uhppote.GetCardsRequest))
		},
	},

	0x5a: &handler{
		func() messages.Request { return new(uhppote.GetCardByIdRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetCardById(rq.(*uhppote.GetCardByIdRequest))
		},
	},

	0x5c: &handler{
		func() messages.Request { return new(uhppote.GetCardByIndexRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetCardByIndex(rq.(*uhppote.GetCardByIndexRequest))
		},
	},

	0x80: &handler{
		func() messages.Request { return new(uhppote.SetDoorDelayRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.SetDoorDelay(rq.(*uhppote.SetDoorDelayRequest))
		},
	},

	0x82: &handler{
		func() messages.Request { return new(uhppote.GetDoorDelayRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetDoorDelay(rq.(*uhppote.GetDoorDelayRequest))
		},
	},

	0x94: &handler{
		func() messages.Request { return new(uhppote.FindDevicesRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.Find(rq.(*uhppote.FindDevicesRequest))
		},
	},

	0x96: &handler{
		func() messages.Request { return new(uhppote.SetAddressRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.SetAddress(rq.(*uhppote.SetAddressRequest))
		},
	},

	0xb0: &handler{
		func() messages.Request { return new(uhppote.GetEventRequest) },
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetEvent(rq.(*uhppote.GetEventRequest))
		},
	},

	0xb2: &handler{
		nil,
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.SetEventIndex(rq.(*messages.SetEventIndexRequest))
		},
	},

	0xb4: &handler{
		nil,
		func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
			return s.GetEventIndex(rq.(*messages.GetEventIndexRequest))
		},
	},
}

func simulate() {
	interrupt := make(chan int, 1)
	bind, err := net.ResolveUDPAddr("udp", ":60000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP bind address [%v]", err)))
		return
	}

	go run(bind, interrupt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	stop(interrupt)

	os.Exit(1)
}

func stop(interrupt chan int) {
	interrupt <- 42
	<-interrupt
}

func run(bind *net.UDPAddr, interrupt chan int) {
	connection, err := net.ListenUDP("udp", bind)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to bind to UDP socket [%v]", err)))
		return
	}

	defer connection.Close()

	wait := make(chan int, 1)
	go func() {
		err := listenAndServe(connection)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0
	}()

	<-interrupt
	connection.Close()
	<-wait
	interrupt <- 0
}

func listenAndServe(c *net.UDPConn) error {
	for {
		request, remote, err := receive(c)
		if err != nil {
			return err
		}

		handle(c, remote, request)
	}

	return nil
}

func handle(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	if len(bytes) != 64 {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid message length %d", len(bytes))))
		return
	}

	if bytes[0] != 0x17 {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid message type 0x%02x", bytes[0])))
		return
	}

	h := handlers[bytes[1]]
	if h == nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid command 0x%02x", bytes[1])))
		return
	}

	if h.factory == nil {
		request, err := messages.UnmarshalRequest(bytes)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			return
		}

		for _, s := range simulators {
			response, err := h.dispatcher(s, *request)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
			} else if response != nil && !reflect.ValueOf(response).IsNil() {
				send(c, src, response)
			}
		}
		return
	}

	// TODO: remove when all requests moved to RequestUnmarshal
	request := h.factory()
	err := codec.Unmarshal(bytes, request)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response, err := h.dispatcher(s, request)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else if response != nil && !reflect.ValueOf(response).IsNil() {
			send(c, src, response)
		}
	}
}

func receive(c *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	request := make([]byte, 2048)

	N, remote, err := c.ReadFromUDP(request)
	if err != nil {
		return []byte{}, nil, errors.New(fmt.Sprintf("Failed to read from UDP socket [%v]", err))
	}

	if options.debug {
		fmt.Printf(" ... received %v bytes from %v\n%s\n", N, remote, dump(request[0:N], " ...          "))
	}

	return request[:N], remote, nil
}

func send(c *net.UDPConn, dest *net.UDPAddr, response interface{}) {
	message, err := codec.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	N, err := c.WriteTo(message, dest)
	if err != nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err)))
	} else {
		if options.debug {
			fmt.Printf(" ... sent %v bytes to %v\n%s\n", N, dest, dump(message[0:N], " ...          "))
		}
	}
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
