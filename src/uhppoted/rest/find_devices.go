package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"uhppote"
	"uhppote/types"
)

type Device struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
	IpAddress    net.IP             `json:"ip-address"`
	SubnetMask   net.IP             `json:"subnet-mask"`
	Gateway      net.IP             `json:"gateway-address"`
	MacAddress   types.MacAddress   `json:"mac-address"`
	Version      types.Version      `json:"version"`
	Date         types.Date         `json:"date"`
}

type DeviceList struct {
	Devices []Device `json:"devices"`
}

func getDevices(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	list := make([]Device, 0)
	for _, d := range devices {
		list = append(list, Device{
			SerialNumber: d.SerialNumber,
			DeviceType:   "UTO311-L04",
			IpAddress:    d.IpAddress,
			SubnetMask:   d.SubnetMask,
			Gateway:      d.Gateway,
			MacAddress:   d.MacAddress,
			Version:      d.Version,
			Date:         d.Date,
		})
	}

	response := DeviceList{list}
	b, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func getDevice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	matches := regexp.MustCompile("^/uhppote/device/([0-9]+)$").FindStringSubmatch(url)
	deviceId, err := strconv.ParseUint(matches[1], 10, 32)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	for _, d := range devices {
		if d.SerialNumber == types.SerialNumber(deviceId) {
			response := Device{
				SerialNumber: d.SerialNumber,
				DeviceType:   "UTO311-L04",
				IpAddress:    d.IpAddress,
				SubnetMask:   d.SubnetMask,
				Gateway:      d.Gateway,
				MacAddress:   d.MacAddress,
				Version:      d.Version,
				Date:         d.Date,
			}

			b, err := json.Marshal(response)
			if err != nil {
				http.Error(w, "Error generating response", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}
	}

	http.Error(w, fmt.Sprintf("No device with ID '%v'", deviceId), http.StatusNotFound)
}
