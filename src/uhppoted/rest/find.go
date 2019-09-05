package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"uhppote"
	"uhppote/types"
)

type device struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
}

func getDevices(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	list := make([]device, 0)
	for _, d := range devices {
		list = append(list, device{
			SerialNumber: d.SerialNumber,
			DeviceType:   "UTO311-L04",
		})
	}

	response := struct {
		Devices []device `json:"devices"`
	}{
		Devices: list,
	}

	reply(ctx, w, response)
}

func getDevice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)

	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	for _, d := range devices {
		if d.SerialNumber == types.SerialNumber(deviceId) {
			response := struct {
				SerialNumber types.SerialNumber `json:"serial-number"`
				DeviceType   string             `json:"device-type"`
				IpAddress    net.IP             `json:"ip-address"`
				SubnetMask   net.IP             `json:"subnet-mask"`
				Gateway      net.IP             `json:"gateway-address"`
				MacAddress   types.MacAddress   `json:"mac-address"`
				Version      types.Version      `json:"version"`
				Date         types.Date         `json:"date"`
			}{
				SerialNumber: d.SerialNumber,
				DeviceType:   "UTO311-L04",
				IpAddress:    d.IpAddress,
				SubnetMask:   d.SubnetMask,
				Gateway:      d.Gateway,
				MacAddress:   d.MacAddress,
				Version:      d.Version,
				Date:         d.Date,
			}

			reply(ctx, w, response)
			return
		}
	}

	http.Error(w, fmt.Sprintf("No device with ID '%v'", deviceId), http.StatusNotFound)
}
