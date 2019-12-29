package uhppoted

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"uhppote"
	"uhppote/types"
)

type DeviceSummary struct {
	DeviceID   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
}

type DeviceDetail struct {
	Device struct {
		ID     uint32 `json:"id"`
		Detail detail `json:"info"`
	} `json:"device"`
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

type GetDevicesRequest struct {
	DeviceID uint32
}

type GetDevicesResponse struct {
	Devices []DeviceSummary `json:"devices"`
}

func (u *UHPPOTED) GetDevices(ctx context.Context, request GetDevicesRequest) (*GetDevicesResponse, int, error) {
	u.debug(ctx, "get-devices", fmt.Sprintf("request  %v", request))

	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()
	if err != nil {
		return nil, StatusInternalServerError, err
	}

	list := make([]DeviceSummary, 0)
	for _, d := range devices {
		item := DeviceSummary{
			DeviceID:   uint32(d.SerialNumber),
			DeviceType: identify(d.SerialNumber),
		}

		list = append(list, item)
	}

	response := GetDevicesResponse{
		Devices: list,
	}

	u.debug(ctx, "get-devices", fmt.Sprintf("response %v", response))

	return &response, StatusOK, nil
}

func (u *UHPPOTED) GetDevice(ctx context.Context, rq Request) {
	u.debug(ctx, "get-device", rq)

	id, err := rq.DeviceID()
	if err != nil {
		u.warn(ctx, 0, "get-device", err)
		u.oops(ctx, "get-device", "Error retrieving device list (invalid device ID)", StatusBadRequest)
		return
	}

	device, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevice(*id)
	if err != nil {
		u.warn(ctx, *id, "get-device", err)
		u.oops(ctx, "get-device", fmt.Sprintf("Error retrieving device summary for '%d'", *id), StatusInternalServerError)
		return
	}

	if device == nil {
		u.warn(ctx, *id, "get-device", fmt.Errorf("No device with ID '%v'", *id))
		u.oops(ctx, "get-device", fmt.Sprintf("Error retrieving device summary for '%d'", id), StatusNotFound)
		return
	}

	response := DeviceDetail{
		struct {
			ID     uint32 `json:"id"`
			Detail detail `json:"info"`
		}{
			ID: *id,
			Detail: detail{
				SerialNumber: device.SerialNumber,
				DeviceType:   identify(device.SerialNumber),
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

func identify(deviceID types.SerialNumber) string {
	id := strconv.FormatUint(uint64(deviceID), 10)

	if strings.HasPrefix(id, "4") {
		return "UTO311-L04"
	}

	if strings.HasPrefix(id, "3") {
		return "UTO311-L03"
	}

	if strings.HasPrefix(id, "2") {
		return "UTO311-L02"
	}

	if strings.HasPrefix(id, "1") {
		return "UTO311-L01"
	}

	return "UTO311-L0x"
}
