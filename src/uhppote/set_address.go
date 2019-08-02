package uhppote

import (
	"errors"
	"fmt"
	"net"
	"uhppote/messages"
	"uhppote/types"
)

func (u *UHPPOTE) SetAddress(serialNumber uint32, address, mask, gateway net.IP) (*types.Result, error) {
	if address.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid IP address: %v", address))
	}

	if mask.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid subnet mask: %v", mask))
	}

	if gateway.To4() == nil {
		return nil, errors.New(fmt.Sprintf("Invalid gateway address: %v", gateway))
	}

	request := messages.SetAddressRequest{
		SerialNumber: types.SerialNumber(serialNumber),
		Address:      address,
		Mask:         mask,
		Gateway:      gateway,
		MagicNumber:  0x55aaaa55,
	}

	// UTC0311-L04 doesn't seem to send a response. The reported remote IP address doesn't change on subsequent commands
	// (both internally and onl Wireshark) but the UTC0311-L04 only replies to ping's on the new IP address. Wireshark
	// reports a 'Gratuitous ARP request' which looks correct after a set-address. Might be something to do with the
	// TPLink or OSX ARP implementation.
	if err := u.Execute(serialNumber, request, nil); err != nil {
		return nil, err
	}

	return &types.Result{
		SerialNumber: types.SerialNumber(serialNumber),
		Succeeded:    true,
	}, nil
}
