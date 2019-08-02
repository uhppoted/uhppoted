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
	"uhppote-simulator/rest"
	"uhppote-simulator/simulator"
	"uhppote-simulator/simulator/entities"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/messages"
)

type handler func(*simulator.Simulator, messages.Request) (messages.Response, error)

var handlers = map[byte]handler{
	0x20: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetStatus(rq.(*messages.GetStatusRequest))
	},

	0x30: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.SetTime(rq.(*messages.SetTimeRequest))
	},

	0x32: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetTime(rq.(*messages.GetTimeRequest))
	},

	0x40: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.OpenDoor(rq.(*messages.OpenDoorRequest))
	},

	0x50: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.PutCard(rq.(*messages.PutCardRequest))
	},

	0x52: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.DeleteCard(rq.(*messages.DeleteCardRequest))
	},

	0x54: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.DeleteCards(rq.(*messages.DeleteCardsRequest))
	},

	0x58: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetCards(rq.(*messages.GetCardsRequest))
	},

	0x5a: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetCardById(rq.(*messages.GetCardByIdRequest))
	},

	0x5c: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetCardByIndex(rq.(*messages.GetCardByIndexRequest))
	},

	0x80: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.SetDoorDelay(rq.(*messages.SetDoorDelayRequest))
	},

	0x82: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetDoorDelay(rq.(*messages.GetDoorDelayRequest))
	},

	0x90: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.SetListener(rq.(*messages.SetListenerRequest))
	},

	0x92: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetListener(rq.(*messages.GetListenerRequest))
	},

	0x94: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.Find(rq.(*messages.FindDevicesRequest))
	},

	0x96: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.SetAddress(rq.(*messages.SetAddressRequest))
	},

	0xb0: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetEvent(rq.(*messages.GetEventRequest))
	},

	0xb2: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.SetEventIndex(rq.(*messages.SetEventIndexRequest))
	},

	0xb4: func(s *simulator.Simulator, rq messages.Request) (messages.Response, error) {
		return s.GetEventIndex(rq.(*messages.GetEventIndexRequest))
	},
}

func simulate() {
	bind, err := net.ResolveUDPAddr("udp", ":60000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP bind address [%v]", err)))
		return
	}

	connection, err := net.ListenUDP("udp", bind)
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to bind to UDP socket [%v]", err)))
		return
	}

	defer connection.Close()

	wait := make(chan int, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go run(connection, wait)

	<-interrupt
	connection.Close()
	<-wait

	os.Exit(1)
}

func run(connection *net.UDPConn, wait chan int) {
	queue := make(chan entities.Message, 8)

	for _, s := range simulators {
		s.TxQueue = queue
	}

	go func() {
		err := listenAndServe(connection, queue)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		wait <- 0
	}()

	go func() {
		for {
			msg := <-queue
			send(connection, msg.Destination, msg.Message)
		}
	}()

	go func() {
		rest.Run(simulators)
	}()

}

func listenAndServe(c *net.UDPConn, queue chan entities.Message) error {
	for {
		request, remote, err := receive(c)
		if err != nil {
			return err
		}

		handle(c, remote, request, queue)
	}

	return nil
}

func handle(c *net.UDPConn, src *net.UDPAddr, bytes []byte, queue chan entities.Message) {
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

	request, err := messages.UnmarshalRequest(bytes)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	for _, s := range simulators {
		response, err := h(s, *request)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else if response != nil && !reflect.ValueOf(response).IsNil() {
			queue <- entities.Message{src, response}
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
	msg, err := codec.Marshal(response)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	N, err := c.WriteTo(msg, dest)
	if err != nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err)))
	} else {
		if options.debug {
			fmt.Printf(" ... sent %v bytes to %v\n%s\n", N, dest, dump(msg[0:N], " ...          "))
		}
	}
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}
