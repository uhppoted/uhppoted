package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"uhppote"
	"uhppote/types"
)

type device struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
}

func getDevices(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	debug(ctx, 0, "get-devices", r)

	u := ctx.Value("uhppote").(*uhppote.UHPPOTE)
	wg := sync.WaitGroup{}
	list := sync.Map{}

	for id, _ := range u.Devices {
		deviceID := id
		wg.Add(1)
		go func() {
			defer wg.Done()
			if device, err := u.FindDevice(deviceID); err != nil {
				warn(ctx, deviceID, "get-devices", err)
			} else if device != nil {
				list.Store(uint32(device.SerialNumber), device)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if devices, err := u.FindDevices(); err != nil {
			warn(ctx, 0, "get-devices", err)
		} else {
			for _, d := range devices {
				list.Store(uint32(d.SerialNumber), d)
			}
		}
	}()

	wg.Wait()

	devices := make([]device, 0)
	list.Range(func(key, value interface{}) bool {
		if d, ok := value.(types.Device); ok {
			devices = append(devices, device{
				SerialNumber: d.SerialNumber,
				DeviceType:   identify(d.SerialNumber),
			})
		}

		if d, ok := value.(*types.Device); ok {
			devices = append(devices, device{
				SerialNumber: d.SerialNumber,
				DeviceType:   identify(d.SerialNumber),
			})
		}

		return true
	})

	response := struct {
		Devices []device `json:"devices"`
	}{
		Devices: devices,
	}

	reply(ctx, w, response)
}

func getDevice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)

	device, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevice(deviceID)

	if err != nil {
		warn(ctx, deviceID, "get-device", err)
		http.Error(w, "Error retrieving device list", http.StatusInternalServerError)
		return
	}

	if device == nil {
		http.Error(w, fmt.Sprintf("No device with ID '%v'", deviceID), http.StatusNotFound)
		return
	}

	response := struct {
		SerialNumber types.SerialNumber `json:"serial-number"`
		DeviceType   string             `json:"device-type"`
		IPAddress    net.IP             `json:"ip-address"`
		SubnetMask   net.IP             `json:"subnet-mask"`
		Gateway      net.IP             `json:"gateway-address"`
		MacAddress   types.MacAddress   `json:"mac-address"`
		Version      types.Version      `json:"version"`
		Date         types.Date         `json:"date"`
	}{
		SerialNumber: device.SerialNumber,
		DeviceType:   identify(device.SerialNumber),
		IPAddress:    device.IpAddress,
		SubnetMask:   device.SubnetMask,
		Gateway:      device.Gateway,
		MacAddress:   device.MacAddress,
		Version:      device.Version,
		Date:         device.Date,
	}

	reply(ctx, w, response)
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
