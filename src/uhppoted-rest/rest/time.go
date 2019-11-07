package rest

import (
	"context"
	"encoding/json"
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
		warn(ctx, deviceId, "get-time", err)
		http.Error(w, "Error retrieving device time", http.StatusInternalServerError)
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
		warn(ctx, deviceId, "set-time", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		DateTime types.DateTime `json:"datetime"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		warn(ctx, deviceId, "set-time", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetTime(deviceId, time.Time(body.DateTime))
	if err != nil {
		warn(ctx, deviceId, "set-time", err)
		http.Error(w, "Error setting device time", http.StatusInternalServerError)
		return
	}

	response := struct {
		DateTime types.DateTime `json:"datetime"`
	}{
		DateTime: result.DateTime,
	}

	reply(ctx, w, response)
}
