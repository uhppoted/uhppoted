package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uhppoted/uhppoted/src/uhppote"
	"io/ioutil"
	"net/http"
)

func getDoorDelay(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceID, door)
	if err != nil {
		warn(ctx, deviceID, "get-door-delay", err)
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
	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, deviceID, "set-door-delay", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		Delay uint8 `json:"delay"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		warn(ctx, deviceID, "set-door-delay", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceID, door)
	if err != nil {
		warn(ctx, deviceID, "set-door-delay", err)
		http.Error(w, "Error setting door delay", http.StatusInternalServerError)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(deviceID, door, state.ControlState, body.Delay)
	if err != nil {
		warn(ctx, deviceID, "set-door-delay", err)
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

func getDoorControl(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceID, door)
	if err != nil {
		warn(ctx, deviceID, "get-door-control", err)
		http.Error(w, "Error retrieving door control", http.StatusInternalServerError)
		return
	}

	response := struct {
		ControlState string `json:"control"`
	}{
		ControlState: lookup[result.ControlState],
	}

	reply(ctx, w, response)
}

func setDoorControl(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	states := map[string]uint8{
		"normally open":   1,
		"normally closed": 2,
		"controlled":      3,
	}

	lookup := map[uint8]string{
		1: "normally open",
		2: "normally closed",
		3: "controlled",
	}

	deviceID := ctx.Value("device-id").(uint32)
	door := ctx.Value("door").(uint8)

	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		warn(ctx, deviceID, "set-door-control", err)
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	body := struct {
		ControlState string `json:"control"`
	}{}

	err = json.Unmarshal(blob, &body)
	if err != nil {
		warn(ctx, deviceID, "set-door-control", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	} else if _, ok := states[body.ControlState]; !ok {
		warn(ctx, deviceID, "set-door-control", fmt.Errorf("Invalid request control value: '%s'", body.ControlState))
		http.Error(w, fmt.Sprintf("Invalid request value: '%s'", body.ControlState), http.StatusBadRequest)
		return
	}

	state, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetDoorControlState(deviceID, door)
	if err != nil {
		warn(ctx, deviceID, "set-door-control", err)
		http.Error(w, "Error setting door control", http.StatusInternalServerError)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).SetDoorControlState(deviceID, door, states[body.ControlState], state.Delay)
	if err != nil {
		warn(ctx, deviceID, "set-door-control", err)
		http.Error(w, "Error setting door control", http.StatusInternalServerError)
		return
	}

	response := struct {
		ControlState string `json:"control"`
	}{
		ControlState: lookup[result.ControlState],
	}

	reply(ctx, w, response)
}
