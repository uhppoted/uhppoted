package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"uhppote"
	"uhppote/types"
)

func getTime(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetTime(deviceId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving device time: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: result.DateTime,
	}

	reply(ctx, w, response)
}

func setTime(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId := ctx.Value("device-id").(uint32)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		DateTime types.DateTime `json:"datetime"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetTime(deviceId, time.Time(body.DateTime))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error setting device time: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: result.DateTime,
	}

	reply(ctx, w, response)
}
