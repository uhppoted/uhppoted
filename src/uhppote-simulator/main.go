package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"uhppote-simulator/simulator"
)

type addr struct {
	address *net.UDPAddr
}

var VERSION = "v0.00.0"

var options = struct {
	dir   string
	debug bool
}{
	dir:   "./devices",
	debug: false,
}

var simulators = []*simulator.Simulator{}

var handlers = map[byte]func(*net.UDPConn, *net.UDPAddr, []byte){
	0x94: find,
	0x5a: getCardById,
}

func main() {
	flag.StringVar(&options.dir, "devices", "./devices", "Specifies the simulation device directory")
	flag.BoolVar(&options.debug, "debug", false, "Displays simulator activity")
	flag.Parse()

	interrupt := make(chan int, 1)
	bind, err := net.ResolveUDPAddr("udp", ":60000")
	if err != nil {
		fmt.Printf("%v\n", errors.New(fmt.Sprintf("Failed to resolve UDP bind address [%v]", err)))
		return
	}

	simulators = load(options.dir)

	go run(bind, interrupt)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	stop(interrupt)

	os.Exit(1)
}

func load(dir string) []*simulator.Simulator {
	if options.debug {
		fmt.Printf("   ... loading devices from '%s'\n", dir)
	}

	simulators := []*simulator.Simulator{}
	files, err := filepath.Glob(path.Join(dir, "*.gz"))
	if err == nil {
		for _, f := range files {
			if options.debug {
				fmt.Printf("   ... loading device from '%s'\n", f)
			}

			s, err := simulator.LoadGZ(f)
			if err != nil {
				fmt.Printf("WARN: could not load device from file '%s': %v\n", f, err)
			} else {
				simulators = append(simulators, s)
			}
		}
	}

	return simulators
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
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid message type %02X", bytes[0])))
		return
	}

	h := handlers[bytes[1]]
	if h == nil {
		fmt.Printf("ERROR: %v\n", errors.New(fmt.Sprintf("Invalid command %02X", bytes[1])))
		return
	}

	h(c, src, bytes)
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

func send(c *net.UDPConn, dest *net.UDPAddr, message []byte) error {
	N, err := c.WriteTo(message, dest)

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to write to UDP socket [%v]", err))
	}

	if options.debug {
		fmt.Printf(" ... sent %v bytes to %v\n%s\n", N, dest, dump(message[0:N], " ...          "))
	}

	return nil
}

func find(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	for _, s := range simulators {
		reply, err := s.Find(bytes)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else if err = send(c, src, reply); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func getCardById(c *net.UDPConn, src *net.UDPAddr, bytes []byte) {
	for _, s := range simulators {
		reply, err := s.GetCardById(bytes)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		} else if err = send(c, src, reply); err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
}

func usage() error {
	fmt.Println("Usage: uhppote-simulator [options]")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --help     Displays this message")
	fmt.Println("               For help on a specific command use 'uhppote-cli help <command>'")
	fmt.Println("    --version  Displays the current version of uhppote-simulator")
	fmt.Println("    --debug    Displays a trace of simulator activity")
	fmt.Println("    --devices  The directory containing the simulator device files")
	fmt.Println()

	return nil
}

func version() error {
	fmt.Printf("%v\n", VERSION)

	return nil
}

func dump(m []byte, prefix string) string {
	regex := regexp.MustCompile("(?m)^(.*)")

	return fmt.Sprintf("%s", regex.ReplaceAllString(hex.Dump(m), prefix+"$1"))
}

func (b *addr) String() string {
	return b.address.String()
}

func (b *addr) Set(s string) error {
	address, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return err
	}

	b.address = address

	return nil
}
