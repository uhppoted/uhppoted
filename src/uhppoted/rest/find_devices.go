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
	URI        string `json:"uri"`
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
			URI:        fmt.Sprintf("/uhppote/device/%d", d.SerialNumber),
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
