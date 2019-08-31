package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"uhppote"
)

type Device struct {
	DeviceId   uint32 `json:"device-id"`
	DeviceType string `json:"device-type"`
	IpAddress  string `json:"ip-address"`
	SubnetMask string `json:"subnet-mask"`
	Gateway    string `json:"gateway-address"`
	MacAddress string `json:"mac-address"`
	Version    string `json:"version"`
	Date       string `json:"date"`
}

type DeviceList struct {
	Devices []Device `json:"devices"`
}

func GetDevices(u *uhppote.UHPPOTE, w http.ResponseWriter, r *http.Request) {
	devices, err := u.FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	list := make([]Device, 0)
	for _, d := range devices {
		list = append(list, Device{
			DeviceId:   uint32(d.SerialNumber),
			DeviceType: "UTO311-L04",
			IpAddress:  d.IpAddress.String(),
			SubnetMask: d.SubnetMask.String(),
			Gateway:    d.Gateway.String(),
			MacAddress: d.MacAddress.String(),
			Version:    fmt.Sprintf("%04x", d.Version),
			Date:       d.Date.String(),
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

func GetDevice(deviceId uint32, u *uhppote.UHPPOTE, w http.ResponseWriter, r *http.Request) {
	devices, err := u.FindDevices()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device list: %v", err), http.StatusInternalServerError)
		return
	}

	for _, d := range devices {
		if uint32(d.SerialNumber) == deviceId {
			response := Device{
				DeviceId:   uint32(d.SerialNumber),
				DeviceType: "UTO311-L04",
				IpAddress:  d.IpAddress.String(),
				SubnetMask: d.SubnetMask.String(),
				Gateway:    d.Gateway.String(),
				MacAddress: d.MacAddress.String(),
				Version:    fmt.Sprintf("%04x", d.Version),
				Date:       d.Date.String(),
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
