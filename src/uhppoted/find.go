package uhppoted

import (
	"context"
	"uhppote"
	"uhppote/types"
)

type device struct {
	SerialNumber types.SerialNumber `json:"serial-number"`
	DeviceType   string             `json:"device-type"`
}

func (u *UHPPOTED) GetDevices(ctx context.Context, request Request) {
	u.debug(ctx, 0, "get-devices", request)

	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()
	if err != nil {
		u.warn(ctx, 0, "get-devices", err)
		u.oops(ctx, "get-devices", "Error retrieving device list", StatusInternalServerError)
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

	u.reply(ctx, response)
}

//func getDevice(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//	deviceId := ctx.Value("device-id").(uint32)
//
//	devices, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).FindDevices()
//
//	if err != nil {
//		warn(ctx, deviceId, "get-device", err)
//		http.Error(w, "Error retrieving device list", http.StatusInternalServerError)
//		return
//	}
//
//	for _, d := range devices {
//		if d.SerialNumber == types.SerialNumber(deviceId) {
//			response := struct {
//				SerialNumber types.SerialNumber `json:"serial-number"`
//				DeviceType   string             `json:"device-type"`
//				IpAddress    net.IP             `json:"ip-address"`
//				SubnetMask   net.IP             `json:"subnet-mask"`
//				Gateway      net.IP             `json:"gateway-address"`
//				MacAddress   types.MacAddress   `json:"mac-address"`
//				Version      types.Version      `json:"version"`
//				Date         types.Date         `json:"date"`
//			}{
//				SerialNumber: d.SerialNumber,
//				DeviceType:   "UTO311-L04",
//				IpAddress:    d.IpAddress,
//				SubnetMask:   d.SubnetMask,
//				Gateway:      d.Gateway,
//				MacAddress:   d.MacAddress,
//				Version:      d.Version,
//				Date:         d.Date,
//			}
//
//			reply(ctx, w, response)
//			return
//		}
//	}
//
//	http.Error(w, fmt.Sprintf("No device with ID '%v'", deviceId), http.StatusNotFound)
//}