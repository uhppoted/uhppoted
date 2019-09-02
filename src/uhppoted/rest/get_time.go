package rest

import (
	"context"
	"fmt"
	"net/http"
	"uhppote"
	"uhppote/types"
)

func getTime(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	deviceId, err := parse(r)
	if err != nil {
		http.Error(w, "Error reading request", http.StatusInternalServerError)
		return
	}

	result, err := ctx.Value("uhppote").(*uhppote.UHPPOTE).GetTime(uint32(deviceId))
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
