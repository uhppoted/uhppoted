package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"uhppote"
)

func getDoorDelay(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorDelay(deviceId, door)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving door delay: %v", err), http.StatusInternalServerError)
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
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	fmt.Println("DEBUG ", string(blob))
	body := struct {
		Delay uint8 `json:"delay"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorDelay(deviceId, door, body.Delay)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error setting door delay: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Delay uint8 `json:"delay"`
	}{
		Delay: result.Delay,
	}

	reply(ctx, w, response)
}
