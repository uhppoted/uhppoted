package types

import (
	"fmt"
	"net"
)

type Version uint16

type Device struct {
	SerialNumber uint32
	IpAddress    net.IP
	SubnetMask   net.IP
	Gateway      net.IP
	MacAddress   net.HardwareAddr
	Version      Version
	Date         Date
}

func (device *Device) String() string {
	return fmt.Sprintf("%v %v %v %v %v %v %v",
		device.SerialNumber,
		device.IpAddress,
		device.SubnetMask,
		device.Gateway,
		device.MacAddress,
		device.Version,
		device.Date)
}
