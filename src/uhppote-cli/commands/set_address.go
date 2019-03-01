package commands

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"uhppote"
)

type SetAddressCommand struct {
	SerialNumber uint32
	Address      net.IP
	Mask         net.IP
	Gateway      net.IP
}

func NewSetAddressCommand() (*SetAddressCommand, error) {
	serialNumber, err := getUint32(1, "Missing serial number", "Invalid serial number: %v")
	if err != nil {
		return nil, err
	}

	if len(flag.Args()) < 3 {
		return nil, errors.New("Missing IP address")
	}

	address := net.ParseIP(flag.Arg(2))

	if address == nil || address.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP address: %v", flag.Arg(2)))
	}

	mask := net.IPv4(255, 255, 255, 0)
	if len(flag.Args()) > 3 {
		mask = net.ParseIP(flag.Arg(3))

		if mask == nil || mask.To4() == nil {
			mask = net.IPv4(255, 255, 255, 0)
		}
	}

	gateway := net.IPv4(0, 0, 0, 0)

	if len(flag.Args()) > 4 {
		gateway = net.ParseIP(flag.Arg(3))
		if gateway == nil || gateway.To4() == nil {
			gateway = net.IPv4(0, 0, 0, 0)
		}
	}

	return &SetAddressCommand{serialNumber, address, mask, gateway}, nil
}

func (c *SetAddressCommand) Execute(u *uhppote.UHPPOTE) error {
	err := u.SetAddress(c.SerialNumber, c.Address, c.Mask, c.Gateway)

	return err
}
