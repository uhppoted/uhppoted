package uhppoted

import (
	"context"
	"fmt"
	"net"
	"uhppote"
	"uhppote/types"
)

type DeviceSummary struct {
	Device struct {
		ID      uint32  `json:"id"`
		Summary summary `json:"info"`
	} `json:"device"`
}

type DeviceDetail struct {
	Device struct {
		ID     uint32 `json:"id"`
		Detail detail `json:"info"`
	} `json:"device"`
}

type summary struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
}

type detail struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
	IpAddress    net.IP             `json:"ip-address"`
	SubnetMask   net.IP             `json:"subnet-mask"`
	Gateway      net.IP             `json:"gateway-address"`
	MacAddress   types.MacAddress   `json:"mac-address"`
	Version      types.Version      `json:"version"`
	Date         types.Date         `json:"date"`
}

func (u *UHPPOTED) GetDevices(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-devices", rq)

	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()
	if err != nil {
		u.warn(ctx, 0, "get-devices", err)
		u.oops(ctx, "get-devices", "Error retrieving device list", StatusInternalServerError)
		return
	}

	list := make([]DeviceSummary, 0)
	for _, d := range devices {
		item := DeviceSummary{
			struct {
				ID      uint32  `json:"id"`
				Summary summary `json:"info"`
			}{
				ID: uint32(d.SerialNumber),
				Summary: summary{
					SerialNumber: d.SerialNumber,
					DeviceType:   "UTO311-L04",
				},
			},
		}

		list = append(list, item)
	}

	response := struct {
		Devices []DeviceSummary `json:"devices"`
	}{
		Devices: list,
	}

	u.reply(ctx, response)
}

func (u *UHPPOTED) GetDevice(ctx context.Context, rq Request) {
	u.debug(ctx, 0, "get-device", rq)

	id, err := rq.DeviceId()
	if err != nil {
		u.warn(ctx, id, "get-device", err)
		u.oops(ctx, "get-device", "Error retrieving device list (invalid device ID)", StatusBadRequest)
		return
	}

	device, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevice(id)
	if err != nil {
		u.warn(ctx, id, "get-device", err)
		u.oops(ctx, "get-device", fmt.Sprintf("Error retrieving device summary for '%d'", id), StatusInternalServerError)
		return
	}

	if device == nil {
		u.warn(ctx, id, "get-device", fmt.Errorf("No device with ID '%v'", id))
		u.oops(ctx, "get-device", fmt.Sprintf("Error retrieving device summary for '%d'", id), StatusNotFound)
		return
	}

	response := DeviceDetail{
		struct {
			ID     uint32 `json:"id"`
			Detail detail `json:"info"`
		}{
			ID: id,
			Detail: detail{
				SerialNumber: device.SerialNumber,
				DeviceType:   "UTO311-L04",
				IpAddress:    device.IpAddress,
				SubnetMask:   device.SubnetMask,
				Gateway:      device.Gateway,
				MacAddress:   device.MacAddress,
				Version:      device.Version,
				Date:         device.Date,
			},
		},
	}

	u.reply(ctx, response)
}
