package rest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"uhppote"
)

func getDoorDelay(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceId, door)
	if err != nil {
		warn(ctx, deviceId, "get-door-delay", err)
		http.Error(w, "Error retrieving door delay", http.StatusInternalServerError)
		return
	}

	response := struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: result.Delay,
	}

	reply(ctx, w, response)
}

func setDoorDelay(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, deviceId, "set-door-delay", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		Delay uint8 `json:"delay"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		warn(ctx, deviceId, "set-door-delay", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(deviceId, door, 3, body.Delay)
	if err != nil {
		warn(ctx, deviceId, "set-door-delay", err)
		http.Error(w, "Error setting door delay", http.StatusInternalServerError)
		return
	}

	response := struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: result.Delay,
	}

	reply(ctx, w, response)
}
