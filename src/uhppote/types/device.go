package types

import (
	"fmt"
	"net"
)

type Version uint16

type Device struct {
	SerialNumber SerialNumber
	IpAddress    net.IP
	SubnetMask   net.IP
	Gateway      net.IP
	MacAddress   net.HardwareAddr
	Version      Version
	Date         Date
}

func (device *Device) String() string {
	return fmt.Sprintf("%s %v %v %v %v %04X %s",
		device.SerialNumber,
		device.IpAddress,
		device.SubnetMask,
		device.Gateway,
		device.MacAddress,
		device.Version,
		device.Date)
}
